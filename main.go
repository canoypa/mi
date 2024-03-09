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

	"github.com/canoypa/mi/auth"
	"github.com/canoypa/mi/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type RequestBody struct {
	I              string   `json:"i"`
	Text           string   `json:"text"`
	Visibility     string   `json:"visibility,omitempty"` // public, home, followers, specified
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
	flagPublic       bool
	flagHomeTimeline bool
	flagFollowers    bool
	flagDirect       []string

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

		hostname := viper.GetString("default.hostname")
		token := viper.GetString("default.token")
		if hostname == "" || token == "" {
			fmt.Println("It seems like it's being executed for the first time.")
			fmt.Println("To use this tool, you must set the hostname and access token.")
			confirmInitialize := utils.Confirm("Would you like to set it now?", true)

			if !confirmInitialize {
				os.Exit(0)
			}

			initialize(cmd)
		}

		text := ""
		if len(args) == 0 {
			text = utils.Multiline(getRandomPlaceholder())
		} else {
			text = strings.Join(args, " ")
		}

		// 空でなければ投稿
		if len(text) > 0 {
			post(text)
		}
	},
}

func post(text string) {
	hostname := viper.GetString("default.hostname")
	token := viper.GetString("default.token")

	requestBody := RequestBody{
		I:    token,
		Text: text,
	}

	if len(flagDirect) > 0 {
		requestBody.Visibility = "specified"
		requestBody.VisibleUserIds = flagDirect
	} else if flagFollowers {
		requestBody.Visibility = "followers"
	} else if flagHomeTimeline {
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

func miAuth(hostname string) string {
	sessionId := auth.NewSessionId()
	authConfig := auth.AuthConfig{
		Name:       "mi",
		Permission: []string{"write:notes"},
	}

	authUrl := auth.NewAuthUrl(hostname, sessionId, authConfig)

	fmt.Println("Please access the following URL and authenticate.")
	fmt.Println(authUrl.String())

	utils.OpenUrl(authUrl)
	utils.Input("Press Enter after authentication.") // only for waiting

	res, err := auth.FetchToken(hostname, sessionId)
	cobra.CheckErr(err)

	return res.Token
}

func initialize(cmd *cobra.Command) {
	fmt.Println("Enter the hostname you wish to use. For example, \"misskey.io\".")
	hostname := utils.Input("Hostname:")

	fmt.Println("Chose the authentication method.")
	authMethod := utils.Select("Authentication method:", []string{"MiAuth", "Access Token"})

	var token string
	if authMethod == "MiAuth" {
		token = miAuth(hostname)
	} else if authMethod == "Access Token" {
		fmt.Println("Enter the access token. \"Compose and delete notes\" permission is required.")
		token = utils.Input("Access Token:")
	}

	viper.Set("default.hostname", hostname)
	viper.Set("default.token", token)

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
  -t, --timeline       Publish Note to home timeline
  -f, --followers      Publish Note to followers
  -d, --direct string  Publish Note to specified users
  -l, --local-only     Only sent note to local
  -w, --cw string      Set contents warning

      --init           Set the host and access token

  -h, --help           Help for mi

Examples:
  $ mi Hello world!
  $ mi --cw Read? It's nsfw!
  $ mi --direct "misskey,misskey@example.com" Hello Misskey!
  $ mi --set visibility=public --set local-only=true
`)

	rootCmd.PersistentFlags().BoolVarP(&flagPublic, "public", "p", true, "Publish Note to all users (default)")
	rootCmd.PersistentFlags().BoolVarP(&flagHomeTimeline, "timeline", "t", false, "Publish Note to home timeline")
	rootCmd.PersistentFlags().BoolVarP(&flagFollowers, "followers", "f", false, "Publish Note to followers")
	rootCmd.PersistentFlags().StringSliceVarP(&flagDirect, "direct", "d", []string{}, "Publish Note to specified users")
	rootCmd.MarkFlagsMutuallyExclusive("public", "timeline", "followers", "direct")

	rootCmd.PersistentFlags().BoolVarP(&flagLocalOnly, "local-only", "l", false, "Publish Note only to local")
	rootCmd.PersistentFlags().StringVarP(&flagCw, "cw", "w", "", "Set contents warning")

	rootCmd.PersistentFlags().BoolVar(&flagInit, "init", false, "Set the host and access token")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
