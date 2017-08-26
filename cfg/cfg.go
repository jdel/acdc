package cfg

import (
	"encoding/base64"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"

	"github.com/jdel/acdc/stc"
	"github.com/jdel/acdc/util"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"gopkg.in/yaml.v2"
)

// Version is the current application version.
// This variable is populated when building the binary with:
// -ldflags "-X cfg.Version=${VERSION}"
var Version,
	appHome string

// DockerComposePath is the path docker-compose is available at
var DockerComposePath string

var logConfig = log.WithFields(log.Fields{
	"module": "config",
})

// InitConfig loads the config file according
// to cfgFile and homeDir flags from cmd/root.go
func InitConfig(cfgFile string, homeDir string) {
	var err error

	// Set logLevel ASAP
	logLevel := parseLogLevel(GetLogLevel())
	log.SetLevel(logLevel)

	// Instantiate appHome ASAP
	appHome = getOrCreateHome(homeDir)

	DockerComposePath, err = exec.LookPath("docker-compose")
	if err != nil {
		logConfig.Fatal("Please install docker-compose on your system")
	}

	logConfig.Info("Docker-compose path: ", DockerComposePath)

	if cfgFile != "" {
		// Use config file from the flag if present
		viper.SetConfigFile(cfgFile)
	} else {
		// Otherwise use home flag if present
		if homeDir != "" {
			viper.AddConfigPath(homeDir)
			viper.SetConfigName("acdc")
		} else {
			//Search config in home directory with name ".config" (without extension).
			viper.AddConfigPath("./acdc")
			viper.AddConfigPath(appHome)
			viper.SetConfigName("acdc")
		}
	}

	// Read the config
	if err := viper.ReadInConfig(); err != nil {
		e, ok := err.(viper.ConfigParseError)
		if ok {
			logConfig.Error(e)
		}

		logConfig.Info("No config file used, writing acdc.yml with default values")
		if masterKey, err := util.GenerateRandomString(16); err != nil {
			logConfig.Error("Could not generate key", err)
			viper.Set("acdc.master-key", "pleasechangeme")
		} else {
			viper.Set("acdc.master-key", masterKey)
		}

		settings, _ := yaml.Marshal(viper.AllSettings())
		if err := ioutil.WriteFile(filepath.Join(appHome, "acdc.yml"), settings, 0644); err != nil {
			logConfig.Error(err)
		}
	}

	logConfig.Info("Using config file: ", viper.ConfigFileUsed())
	logConfig.Info("Home is: ", appHome)

	// Create directories
	var staticDir string
	if staticDir, err = util.CreateDir(GetStaticDir()); err != nil {
		logConfig.Fatal("Could not create directory", GetStaticDir())
	}
	if _, err := util.CreateDir(GetComposeDir()); err != nil {
		logConfig.Fatal("Could not create directory", GetComposeDir())
	}

	// Write the static.png file, always overwrite
	// TODO: Remove as i will most likely not need any static stuff
	data, _ := base64.StdEncoding.DecodeString(stc.StaticPng)
	ioutil.WriteFile(filepath.Join(staticDir, "static.png"), data, 0664)

	// Watch for changes
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		logConfig.Info("Config file changed: ", e.Name)
		logLevel := parseLogLevel(GetLogLevel())
		log.SetLevel(logLevel)
	})
}

func parseLogLevel(level string) log.Level {
	var logLevel log.Level
	var err error
	logConfig.WithField("log-level", level).Debug("Parsing log level")
	if logLevel, err = log.ParseLevel(level); err != nil {
		logLevel = log.ErrorLevel
		logConfig.WithField("log-level", level).Error("Cannot parse log level, setting to Error")
	}
	return logLevel
}

// GetPort returns the port
// default value is "8080"
func GetPort() string {
	return viper.GetString("acdc.port")
}

// GetLogLevel returns the log level.
// default value is "Error"
func GetLogLevel() string {
	return viper.GetString("acdc.log-level")
}

// GetMasterKey returns the master api key
// default value is ""
func GetMasterKey() string {
	return viper.GetString("acdc.master-key")
}

// GetStaticPrefix returns the prefix
// for static URLs
func GetStaticPrefix() string {
	return viper.GetString("acdc.router.static")
}

// GetStaticDir returns the static directory
func GetStaticDir() string {
	dir := viper.GetString("acdc.filesystem.static-dir")
	if match, _ := regexp.MatchString("^/", dir); !match {
		return filepath.Join(appHome, dir)
	}
	return dir
}

// GetComposeDir returns the compose directory
func GetComposeDir() string {
	dir := viper.GetString("acdc.filesystem.compose-dir")
	if match, _ := regexp.MatchString("^/", dir); !match {
		return filepath.Join(appHome, dir)
	}
	return dir
}

// getOrCreateHome returns acdc subdir from
// user's home directory and creates it if required
func getOrCreateHome(appHome string) string {
	var home string
	if appHome != "" {
		home = appHome
	} else {
		usr, err := user.Current()
		if err != nil {
			logConfig.Fatal(err)
		}
		logConfig.Info("Current user: ", usr.Username)
		home = filepath.Join(usr.HomeDir, "/acdc/")
	}

	if !util.FileExists(home) {
		logConfig.Info("Creating home: ", home)
		if err := os.Mkdir(home, 0755); err != nil {
			logConfig.Fatal(err)
		}
	}
	return home
}

func init() {
	// Sets logrus options
	formatter := &log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "06/01/02 15:04:05.000",
	}
	log.SetFormatter(formatter)
	log.SetOutput(os.Stderr)
}
