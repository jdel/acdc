package cmd // import jdel.org/acdc/cmd

import (
	"fmt"

	"jdel.org/acdc/api"
	"github.com/spf13/cobra"
)

var apiKeyListCmd = &cobra.Command{
	Use: "list",
	Aliases: []string{
		"ls",
	},
	Short: "List all api-keys",
	Long:  `List all api-keys`,
	Run: func(cmd *cobra.Command, args []string) {
		keys, _ := api.AllAPIKeys()
		for _, key := range keys {
			fmt.Println(key.Unique, "\t", key.Remote)
		}
	},
}

func init() {
	apiKeyCmd.AddCommand(apiKeyListCmd)
}
