package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/canoypa/mi/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type RequestBody struct {
	I              string   `json:"i"`
	Text           string   `json:"text"`
	Visibility     string   `json:"visibility,omitempty"`
	VisibleUserIds []string `json:"visibleUserIds,omitempty"`
	Cw             string   `json:"cw,omitempty"`
	LocalOnly      bool     `json:"localOnly,omitempty"`
}

type Note struct {
	Id        string
	CreatedAt string
	Text      string
	Cw        string
	// User
	UserId     string
	Visibility string
}
type CreateResponse struct {
	CreatedNote Note
}

var (
	flagPublic    bool
	flagTimeline  bool
	flagFollowers bool
	flagDirect    string

	flagLocalOnly bool
	flagCw        string

	flagInit bool
)

func getRandomPlaceholder() string {
	words := []string{
		"What are you up to",
		"What's happening around you",
		"What's on your mind",
		"What do you want to say",
		"Start writing...",
		"Waiting for you to write...",
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	num := r.Intn(len(words))

	return words[num]
}

var rootCmd = &cobra.Command{
	Use:   "mi <text> [flags]",
	Short: "Misskey CLI",
	Long:  "CLI tool for sending Misskey notes.",
	Run: func(cmd *cobra.Command, args []string) {
		if flagInit {
			initialize(cmd)
			os.Exit(0)
		}

		hostname := viper.GetString("hostname")
		token := viper.GetString("token")
		if hostname == "" || token == "" {
			fmt.Println("It seems like it's being executed for the first time.")
			confirmInitialize := utils.Confirm("Would you like to set the hostname and access token?:", true)

			if confirmInitialize {
				initialize(cmd)
				os.Exit(0)
			}
		}

		text := ""
		if len(args) == 0 {
			text = utils.Multiline(getRandomPlaceholder())
		} else {
			text = args[0]
		}

		// 空でなければ投稿
		if len(text) > 0 {
			post(text)
		}
	},
}

func post(text string) {
	hostname := viper.GetString("hostname")
	token := viper.GetString("token")

	requestBody := RequestBody{
		I:    token,
		Text: text,
	}

	if flagDirect != "" {
		requestBody.Visibility = "specified"

		visibleUserIds := strings.Split(flagDirect, ",")
		requestBody.VisibleUserIds = visibleUserIds
	} else if flagFollowers {
		requestBody.Visibility = "followers"
	} else if flagTimeline {
		requestBody.Visibility = "home"
	}

	if flagCw != "" {
		requestBody.Cw = flagCw
	}

	if flagLocalOnly {
		requestBody.LocalOnly = true
	}

	url := url.URL{
		Scheme: "https",
		Host:   hostname,
		Path:   "api/notes/create",
	}

	bodyJson, err := json.Marshal(requestBody)
	cobra.CheckErr(err)

	req, err := http.NewRequest("POST", url.String(), bytes.NewBuffer(bodyJson))
	cobra.CheckErr(err)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	cobra.CheckErr(err)
	defer res.Body.Close()

	if res.StatusCode != 200 {
		fmt.Println("Unknown error!")
		os.Exit(1)
	}

	b, err := io.ReadAll(res.Body)
	cobra.CheckErr(err)

	var response CreateResponse
	err = json.Unmarshal(b, &response)
	cobra.CheckErr(err)

	fmt.Println("Your note was sent: " + "https://" + hostname + "/notes/" + response.CreatedNote.Id)
}

func initialize(cmd *cobra.Command) {
	fmt.Println("Enter the hostname you wish to use. For example, \"misskey.io\".")
	hostname := utils.Input("Hostname:")
	fmt.Println("Enter the access token. \"Compose and delete notes\" permission is required.")
	token := utils.Input("Access Token:")

	viper.Set("hostname", hostname)
	viper.Set("token", token)

	err := viper.WriteConfig()
	cobra.CheckErr(err)

	fmt.Println("Initialization has been completed!")
}

func initConfig() {
	homePath, err := os.UserHomeDir()
	cobra.CheckErr(err)

	configPath := filepath.Join(homePath, ".mi")
	configName := "hosts"
	configType := "yaml"

	viper.AddConfigPath(configPath)
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)

	// if config not found
	if err := viper.ReadInConfig(); err != nil {
		os.MkdirAll(configPath, 0700)
		viper.WriteConfigAs(filepath.Join(configPath, fmt.Sprintf("%s.%s", configName, configType)))
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.SetHelpTemplate(`CLI tool for sending Misskey note.

Usage:
  mi <text> [flags]

Flags:
  -p, --public         Publish Note to all users (default)
  -h, --home           Publish Note to home timeline
  -f, --followers      Publish Note to followers
  -d, --direct string  Publish Note to specified users
      --local-only     Only sent note to local
      --cw string      Set contents warning

      --init           Set the host and access token

  -h, --help           Help for mi

Examples:
  $ mi Hello world!
  $ mi It's nsfw! --cw Read?
  $ mi Hello Misskey! --direct "@misskey,@misskey@example.com"
  $ mi --set visibility=public --set local-only=true
`)

	rootCmd.PersistentFlags().BoolVarP(&flagPublic, "public", "p", true, "Publish Note to all users (default)")
	rootCmd.PersistentFlags().BoolVarP(&flagTimeline, "timeline", "t", false, "Publish Note to home timeline")
	rootCmd.PersistentFlags().BoolVarP(&flagFollowers, "followers", "f", false, "Publish Note to followers")
	rootCmd.PersistentFlags().StringVarP(&flagDirect, "direct", "d", "", "Publish Note to specified users")
	rootCmd.MarkFlagsMutuallyExclusive("public", "timeline", "followers", "direct")

	rootCmd.PersistentFlags().BoolVar(&flagLocalOnly, "local-only", false, "Do not expand mentions from text")
	rootCmd.PersistentFlags().StringVar(&flagCw, "cw", "", "Set contents warning")

	rootCmd.PersistentFlags().BoolVar(&flagInit, "init", false, "Set the host and access token")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
