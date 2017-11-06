package kanal

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "kanal <subcommand>",
	Short: "Telegram bot",
	Run:   nil,
}

func init() {
	cobra.OnInitialize()
}
