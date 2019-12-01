package cmd // import jdel.org/acdc/cmd

import (
	"fmt"
	"runtime"

	"jdel.org/acdc/cfg"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use: "version",
	Aliases: []string{
		"v",
	},
	Short: "Get the version of acdc",
	Long:  `Get the version of acdc`,
	Run: func(cmd *cobra.Command, args []string) {
		v := struct {
			Version   string `json:"version"`
			GoVersion string `json:"go-version"`
			Os        string `json:"os"`
			Arch      string `json:"arch"`
		}{
			Version:   cfg.Version,
			GoVersion: runtime.Version(),
			Os:        runtime.GOOS,
			Arch:      runtime.GOARCH,
		}

		fmt.Println("Version:", v.Version)
		fmt.Println("Go version:", v.GoVersion)
		fmt.Println("Os:", v.Os)
		fmt.Println("Arch:", v.Arch)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
