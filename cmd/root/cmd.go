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
	Command.SetHelpTemplate(`CLI tool for sending Misskey note.

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
