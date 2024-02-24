package main

import (
	"os"
	"path/filepath"

	"github.com/mitsukaki/lumi/internal/api"
	"github.com/mitsukaki/lumi/internal/db"
	"github.com/spf13/viper"
)

func loadConfig() {
	exePath, _ := os.Executable()
	binaryPath := filepath.Dir(exePath)

	viper.SetEnvPrefix("lumi")
	viper.AutomaticEnv()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// set the default values
	viper.SetDefault("db.use_auth", false)
	viper.SetDefault("db.use_https", false)
	viper.SetDefault("db.host", "127.0.0.1")
	viper.SetDefault("db.port", 5984)
	viper.SetDefault("db.username", "")
	viper.SetDefault("db.password", "")
	
	viper.SetDefault("http.port", 8080)

	viper.SetDefault("static_dir", filepath.Join(binaryPath, "public"))

	// read the config file
	viper.AddConfigPath(binaryPath)
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/lumi/")
	viper.AddConfigPath("$HOME/.lumi")

	viper.ReadInConfig()

	// write the default config if needed
	viper.SafeWriteConfig()
}

func main() {
	// load the config
	loadConfig()

	// create the API instance
	apiConfig := api.APIConfig{
		StaticDir: viper.GetString("static_dir"),
		DBConfig: &db.CouchDBConfig{
			UseAuth:  viper.GetBool("db.use_auth"),
			UseHTTPS: viper.GetBool("db.use_https"),
			Host:     viper.GetString("db.host"),
			Port:     viper.GetInt("db.port"),
			Username: viper.GetString("db.username"),
			Password: viper.GetString("db.password"),
		},
	}

	apiServer, err := api.CreateAPIServer(apiConfig)
	if err != nil {
		panic(err)
	}

	apiServer.Start()
}
