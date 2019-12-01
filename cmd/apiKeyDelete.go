package cmd // import jdel.org/acdc/cmd

import (
	"fmt"

	"jdel.org/acdc/api"
	"github.com/spf13/cobra"
)

var apiKeyAddCmd = &cobra.Command{
	Use: "delete <key-unique>",
	Aliases: []string{
		"rm",
	},
	Args:  cobra.ExactArgs(1),
	Short: "Delete an api-key",
	Long:  `Delete an api-key`,
	Run: func(cmd *cobra.Command, args []string) {
		if key := api.FindKey(args[0]); key == nil {
			logCmd.Fatalf("Cannot find key %s", args[0])
		} else {
			if err := key.Delete(); err != nil {
				logCmd.Fatalf("Cannot delete key %s: %s", key.Unique, err)
			}
			fmt.Println("Deleted key: ", key.Unique)
		}
	},
}

func init() {
	apiKeyCmd.AddCommand(apiKeyAddCmd)
}
