package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/canoypa/mi/cmd/root"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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
}

func main() {
	if err := root.Command.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
