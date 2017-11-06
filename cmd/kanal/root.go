package main

import (
	"github.com/spf13/cobra"
	"gitlab.com/arha/kanal/configuration"
)

var rootCmd = &cobra.Command{
	Use:   "kanal <subcommand>",
	Short: "Telegram bot",
	Run:   nil,
}

func init() {
	cobra.OnInitialize(func() {
		configuration.LoadConfig()
	})
}
