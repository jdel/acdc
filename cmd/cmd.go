package cmd

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/jdel/acdc/cfg"
	"github.com/jdel/acdc/rtr"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

var cfgFile, appHome string

var logCmd = log.WithFields(log.Fields{
	"module": "cmd",
})

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "acdc",
	Short: "A Continuous Docker Compose.",
	Long:  `A Continuous Docker Compose provides a docker-compose REST API and hooks for Slack, Docker Hub, Github and more.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Listening on", cfg.GetPort())
		router := rtr.NewRouter()
		logCmd.Fatal(http.ListenAndServe(":"+cfg.GetPort(), router))
	},
}

// Execute runs the main command that
// serves synology packages
func Execute() {
	RootCmd.Execute()
}

func init() {
	// CMD line args > ENV VARS > Config file
	cobra.OnInitialize(func() { cfg.InitConfig(cfgFile, appHome) })
	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "C", "", "config file (default is $HOME/acdc/config.yml)")
	RootCmd.PersistentFlags().StringVarP(&appHome, "home", "H", "", "acdc home (default is $HOME/acdc/)")
	// Optional flags
	RootCmd.PersistentFlags().IntP("port", "p", 8080, "port to listen to")
	RootCmd.PersistentFlags().String("docker-compose-binary", "", "location of the docker-compose binary, empty for auto-detect")
	RootCmd.PersistentFlags().String("compose-dir", "compose", "compose directory")
	RootCmd.PersistentFlags().String("static-dir", "static", "static directory")
	RootCmd.PersistentFlags().StringP("log-level", "l", "Error", "log level [Error,Warn,Info,Debug]")
	RootCmd.PersistentFlags().String("static", "static", "prefix to serve static images")
	// Bind flags to config
	viper.BindPFlag("acdc.port", RootCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("acdc.filesystem.compose-dir", RootCmd.PersistentFlags().Lookup("compose-dir"))
	viper.BindPFlag("acdc.filesystem.static-dir", RootCmd.PersistentFlags().Lookup("static-dir"))
	viper.BindPFlag("acdc.log-level", RootCmd.PersistentFlags().Lookup("log-level"))
	viper.BindPFlag("acdc.router.static", RootCmd.PersistentFlags().Lookup("static"))
	viper.BindPFlag("docker-compose.binary", RootCmd.PersistentFlags().Lookup("docker-compose-binary"))
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
}
