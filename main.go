package main

import (
	"fmt"
	"os"

	"github.com/canoypa/mi/cmd/root"
	"github.com/canoypa/mi/core/config"
	"github.com/spf13/cobra"
)

func init() {
	cobra.OnInitialize(config.InitConfig)
}

func main() {
	if err := root.Command.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
