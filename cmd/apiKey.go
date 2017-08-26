package cmd

import (
	"github.com/spf13/cobra"
)

var apiKeyCmd = &cobra.Command{
	Use: "api-key",
	Aliases: []string{
		"key",
		"k",
	},
	Short: "Make operations on api-keys",
	Long:  `Make operations on api-keys`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func init() {
	RootCmd.AddCommand(apiKeyCmd)
}
