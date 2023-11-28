package config

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type DatabaseConfiguration struct {
	Username string
	Password string
	Database string
	Host     string
	Port     int
}

func (c DatabaseConfiguration) ConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.Username, c.Password, c.Host, c.Port, c.Database)
}

type ScrapeConfiguration struct {
	Interval time.Duration
}

type Configuration struct {
	Database DatabaseConfiguration
	Scrape   ScrapeConfiguration
}

var Config Configuration

func init() {
	cobra.OnInitialize(loadConfiguration)
}

func loadConfiguration() {
	viper.SetDefault("scrape.interval", 10*time.Hour.Seconds())

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	viper.AutomaticEnv()
	viper.BindEnv("database.username", "RECORDER_DB_USERNAME")
	viper.BindEnv("database.password", "RECORDER_DB_PASSWORD")
	viper.BindEnv("database.database", "RECORDER_DB_DATABASE")
	viper.BindEnv("database.host", "RECORDER_DB_HOST")
	viper.BindEnv("scrape.interval", "RECORDER_SCRAPE_INTERVAL")

	requiredKeys := []string{
		"database.username",
		"database.password",
		"database.database",
		"database.host",
		"database.port",
	}
	for _, key := range requiredKeys {
		if !viper.IsSet(key) {
			panic(key + " is required")
		}
	}

	var config Configuration
	if err := viper.Unmarshal(&config); err != nil {
		panic(err)
	}

	Config = config
}
