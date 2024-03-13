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
	Run: func(cmd *cobra.Command, args []string) {
		if flags.FlagInit {
			err := initialize.Command.Execute()
			cobra.CheckErr(err)
			os.Exit(0)
		}

		err := post.Command.Execute()
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

	post.InitFlags(Command)
	initialize.InitFlags(Command)
}
