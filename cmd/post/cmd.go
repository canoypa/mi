package post

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/canoypa/mi/cmd/initialize"
	"github.com/canoypa/mi/core/flags"
	"github.com/canoypa/mi/misskey"
	"github.com/canoypa/mi/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Command(command *cobra.Command, args []string) error {
	hostname := viper.GetString("default.hostname")
	token := viper.GetString("default.token")

	if hostname == "" || token == "" {
		fmt.Println("It seems like it's being executed for the first time.")
		fmt.Println("To use this tool, you must set the hostname and access token.")
		confirmInitialize := utils.Confirm("Would you like to set it now?", true)

		if !confirmInitialize {
			os.Exit(0)
		}

		err := initialize.Command(command, args)
		cobra.CheckErr(err)
	}

	text := ""
	if len(args) == 0 {
		text = utils.Multiline(getRandomPlaceholder())
	} else {
		text = strings.Join(args, " ")
	}

	if text == "" {
		fmt.Println("The note is empty.")
		os.Exit(0)
	}

	post(text)

	return nil
}

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
