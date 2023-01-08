package scraper

import (
	"fmt"
	"os"
	"path/filepath"
)
import "github.com/spf13/viper"

type config struct {
	SearchResource   []string
	SearchKey        string
	OrderBy          string
	MaxResultPerCall int
	ApiKey           []string
}

const (
	scraperConfigName string = "scraper-dev"
	configDirectory   string = "configs"
)

func GetConf() *config {
	conf := &config{}

	currDirectory, err := os.Getwd()
	if err != nil {
		fmt.Printf("error getting the current directory: %v", err)
	}
	configDirectory := filepath.Join(currDirectory, configDirectory)

	viper.AddConfigPath(configDirectory)
	viper.SetConfigName(scraperConfigName)

	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("%v", err)
	}

	err = viper.Unmarshal(conf)
	if err != nil {
		fmt.Printf("unable to decode into config struct, %v", err)
	}

	return conf
}
