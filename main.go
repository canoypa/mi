package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/canoypa/mi/cmd/initialize"
	"github.com/canoypa/mi/core/flags"
	"github.com/canoypa/mi/misskey"
	"github.com/canoypa/mi/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		if flags.FlagInit {
			err := initialize.Command.Execute()
			cobra.CheckErr(err)
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

			err := initialize.Command.Execute()
			cobra.CheckErr(err)
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

	requestBody := misskey.NotesCreateRequestBody{
		I:    token,
		Text: text,
	}

	if len(flags.FlagDirect) > 0 {
		requestBody.Visibility = "specified"
		requestBody.VisibleUserIds = flags.FlagDirect
	} else if flags.FlagFollowers {
		requestBody.Visibility = "followers"
	} else if flags.FlagHomeTimeline {
		requestBody.Visibility = "home"
	}

	if flags.FlagCw != "" {
		requestBody.Cw = flags.FlagCw
	}

	if flags.FlagLocalOnly {
		requestBody.LocalOnly = true
	}

	res, err := misskey.NotesCreate(hostname, requestBody)
	cobra.CheckErr(err)

	fmt.Println("Your note was sent: " + "https://" + hostname + "/notes/" + res.CreatedNote.Id)
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

	rootCmd.PersistentFlags().BoolVarP(&flags.FlagPublic, "public", "p", true, "Publish Note to all users (default)")
	rootCmd.PersistentFlags().BoolVarP(&flags.FlagHomeTimeline, "timeline", "t", false, "Publish Note to home timeline")
	rootCmd.PersistentFlags().BoolVarP(&flags.FlagFollowers, "followers", "f", false, "Publish Note to followers")
	rootCmd.PersistentFlags().StringSliceVarP(&flags.FlagDirect, "direct", "d", []string{}, "Publish Note to specified users")
	rootCmd.MarkFlagsMutuallyExclusive("public", "timeline", "followers", "direct")

	rootCmd.PersistentFlags().BoolVarP(&flags.FlagLocalOnly, "local-only", "l", false, "Publish Note only to local")
	rootCmd.PersistentFlags().StringVarP(&flags.FlagCw, "cw", "w", "", "Set contents warning")

	initialize.InitFlags(rootCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
