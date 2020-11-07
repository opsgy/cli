package opsgy

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	homedir "github.com/mitchellh/go-homedir"

	"gopkg.in/yaml.v3"
)

type Config struct {
	AccessToken  string    `yaml:"accessToken,omitempty"`
	RefreshToken string    `yaml:"refreshToken,omitempty"`
	TokenExpiry  time.Time `yaml:"tokenExpiry,omitempty"`
	TokenType    string    `yaml:"tokenType,omitempty"`

	ProjectName string `yaml:"projectName,omitempty"`
}

var ConfigFile = ""

func SetConfigFile(configFile string) {
	ConfigFile = configFile
}

func LoadConfig() *Config {
	if ConfigFile == "" {
		ConfigFile = GetDefaultConfigFile()
	}
	var config = Config{}
	yamlFile, err := ioutil.ReadFile(ConfigFile)
	if err != nil {
		return &config
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("Failed to not load %s: %v", ConfigFile, err)
	}
	return &config
}

func SaveConfig(config *Config) error {
	if ConfigFile == "" {
		ConfigFile = GetDefaultConfigFile()
	}
	rawYaml, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(ConfigFile, rawYaml, 0600)
	if err != nil {
		return err
	}
	return nil
}

func GetDefaultConfigFile() string {
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return home + "/.opsgy.yml"
}
