package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/canoypa/mi/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	flagPublic    bool
	flagTimeline  bool
	flagFollowers bool
	flagDirect    string

	flagCw string

	flagNoMentions bool
	flagNoHashtags bool
	flagNoEmoji    bool

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
			confirmInitialize := utils.Confirm("Would you like to set the host and access token?:", true)

			if confirmInitialize {
				initialize(cmd)
				os.Exit(0)
			}
		}
	},
}

func initialize(cmd *cobra.Command) {
	fmt.Println("Enter the hostname you wish to use. For example, \"misskey.io\".")
	hostname := utils.Input("Hostname:")
	fmt.Println("Enter the access token. \"Compose and delete notes\" permission is required.")
	token := utils.Input("Access Token:")

	viper.Set("hostname", hostname)
	viper.Set("token", token)

	if err := viper.WriteConfig(); err != nil {
		cobra.CheckErr(err)
	}

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

	rootCmd.SetHelpTemplate(`CLI tool for sending Misskey notes.

Usage:
  mi <text> [flags]

Flags:
  -p, --public         Publish Note to all users (default)
  -t, --timeline       Publish Note to home timeline
  -f, --followers      Publish Note to followers
  -d, --direct string  Publish Note to specified users
      --cw string      Set contents warning

      --no-mentions    Do not expand mentions from text
      --no-hashtags    Do not expand hashtags from text
      --no-emojis       Do not expand emojis from text

      --init           Set the host and access token

  -h, --help           help for mi

Examples:
  $ mi Hello world!
  $ mi It's nsfw! --cw Read?
  $ mi Hello Misskey! --direct @misskey,@example.com@misskey
`)

	rootCmd.PersistentFlags().BoolVarP(&flagPublic, "public", "p", true, "Publish Note to all users (default)")
	rootCmd.PersistentFlags().BoolVarP(&flagTimeline, "timeline", "t", false, "Publish Note to home timeline")
	rootCmd.PersistentFlags().BoolVarP(&flagFollowers, "followers", "f", false, "Publish Note to followers")
	rootCmd.PersistentFlags().StringVarP(&flagDirect, "direct", "d", "", "Publish Note to specified users")
	rootCmd.MarkFlagsMutuallyExclusive("public", "timeline", "followers", "direct")

	rootCmd.PersistentFlags().StringVar(&flagCw, "cw", "", "Set contents warning")

	rootCmd.PersistentFlags().BoolVar(&flagNoMentions, "no-mentions", false, "Do not expand mentions from text")
	rootCmd.PersistentFlags().BoolVar(&flagNoHashtags, "no-hashtags", false, "Do not expand hashtags from text")
	rootCmd.PersistentFlags().BoolVar(&flagNoEmoji, "no-emojis", false, "Do not expand emojis from text")

	rootCmd.PersistentFlags().BoolVar(&flagInit, "init", false, "Set the host and access token")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
