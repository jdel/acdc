package cmd // import jdel.org/acdc/cmd

import (
	"fmt"

	"jdel.org/acdc/api"
	"github.com/spf13/cobra"
)

var flagRemote,
	flagUnique string

var apiKeyNewCmd = &cobra.Command{
	Use: "new",
	Aliases: []string{
		"add",
		"create",
	},
	Short: "Add a new api-key",
	Long:  `Add a new api-key`,
	Run: func(cmd *cobra.Command, args []string) {
		if key, err := api.NewKey(flagUnique, flagRemote); err != nil {
			logCmd.Fatal("Could not create API Key")
		} else {
			fmt.Println(key.Unique, "\t", key.Remote)
		}
	},
}

func init() {
	apiKeyNewCmd.Flags().StringVarP(&flagUnique, "unique", "u", "", "unique id")
	apiKeyNewCmd.Flags().StringVarP(&flagRemote, "remote", "r", "", "remote git repository")
	apiKeyCmd.AddCommand(apiKeyNewCmd)
}
