package root

import (
	"os"

	"github.com/canoypa/mi/cmd/initialize"
	"github.com/canoypa/mi/cmd/post"
	"github.com/canoypa/mi/core/flags"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "mi <text> [flags]",
	Short: "Misskey CLI",
	Long:  "CLI tool for sending Misskey notes.",
	Example: `$ mi Hello world!
$ mi --cw Read? It's nsfw!
$ mi --direct "misskey,misskey@example.com" Hello Misskey!
$ mi --set visibility=public --set local-only=true`,
	Run: func(command *cobra.Command, args []string) {
		if flags.FlagInit {
			err := initialize.Command(command, args)
			cobra.CheckErr(err)
			os.Exit(0)
		}

		err := post.Command(command, args)
		cobra.CheckErr(err)
	},
}

func init() {
	// post command flags
	Command.PersistentFlags().BoolVarP(&flags.FlagPublic, "public", "p", true, "Publish Note to all users (default)")
	Command.PersistentFlags().BoolVarP(&flags.FlagHomeTimeline, "timeline", "t", false, "Publish Note to home timeline")
	Command.PersistentFlags().BoolVarP(&flags.FlagFollowers, "followers", "f", false, "Publish Note to followers")
	Command.PersistentFlags().StringSliceVarP(&flags.FlagDirect, "direct", "d", []string{}, "Publish Note to specified users")
	Command.MarkFlagsMutuallyExclusive("public", "timeline", "followers", "direct")

	Command.PersistentFlags().BoolVarP(&flags.FlagLocalOnly, "local-only", "l", false, "Publish Note only to local")
	Command.PersistentFlags().StringVarP(&flags.FlagCw, "cw", "w", "", "Set contents warning")

	// initial command flags
	Command.PersistentFlags().BoolVar(&flags.FlagInit, "init", false, "Set the host and access token")
}
