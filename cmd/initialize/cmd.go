package initialize

import (
	"fmt"
	"os"

	"github.com/canoypa/mi/misskey"
	"github.com/canoypa/mi/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Command(command *cobra.Command, args []string) error {
	initialize()

	return nil
}

func miAuth(hostname string) string {
	sessionId := misskey.NewSessionId()
	authConfig := misskey.MiAuthConfig{
		Name:       "Misskey CLI",
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
	if hostname == "" {
		fmt.Println("Please enter the hostname.")
		os.Exit(0)
	}

	fmt.Println("Chose the authentication method.")
	authMethod := utils.Select("Authentication method:", []string{"MiAuth", "Access Token"})

	var token string
	if authMethod == "MiAuth" {
		token = miAuth(hostname)
	} else if authMethod == "Access Token" {
		fmt.Println("Enter the access token. \"Compose and delete notes\" permission is required.")
		token = utils.Input("Access Token:")

		if token == "" {
			fmt.Println("Please enter the access token.")
			os.Exit(0)
		}
	}

	viper.Set("default.hostname", hostname)
	viper.Set("default.token", token)

	err := viper.WriteConfig()
	cobra.CheckErr(err)

	fmt.Println("Initialization has been completed!")
}
