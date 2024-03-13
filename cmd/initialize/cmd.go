package initialize

import (
	"fmt"

	"github.com/canoypa/mi/core/flags"
	"github.com/canoypa/mi/misskey"
	"github.com/canoypa/mi/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Command = &cobra.Command{
	Run: func(command *cobra.Command, args []string) {
		initialize()
	},
}

func InitFlags(command *cobra.Command) {
	command.PersistentFlags().BoolVar(&flags.FlagInit, "init", false, "Set the host and access token")
}

func init() {
	InitFlags(Command)
}

func miAuth(hostname string) string {
	sessionId := misskey.NewSessionId()
	authConfig := misskey.MiAuthConfig{
		Name:       "mi",
		Permission: []string{"write:notes"},
	}

	authUrl := misskey.NewMiAuthUrl(hostname, sessionId, authConfig)

	fmt.Println("Please access the following URL and authenticate.")
	fmt.Println(authUrl.String())

	utils.OpenUrl(authUrl)
	utils.Input("Press Enter after authentication.") // only for waiting

	res, err := misskey.MiAuthCheck(hostname, sessionId)
	cobra.CheckErr(err)

	return res.Token
}

func initialize() {
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
