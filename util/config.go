package util

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"path"
)

type Jiraconf struct {
	Username, Apikey, Orgname, CustomDomain, AccessToken string
}

type Config struct {
	Jira     Jiraconf
	IsDev    bool
	FirstRun bool
}

var (
	CONFIG_DIR      = ".config/zilla"
	CONFIG_FILENAME = "config.toml"
	CACHE_FILENAME  = "cache.json"
	LOG_FILENAME    = "log.txt"
)

func getConfigFileIfExists() (*Config, error) {
	config := Config{}
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, errors.New("couldn't locate home directory, falling back to env based config")
	}
	var conf Config
	configPath := path.Join(home, CONFIG_DIR, CONFIG_FILENAME)
	_, err = toml.Decode(configPath, &conf)
	if err != nil {
		return nil, fmt.Errorf("error parsing config at \"%v\" falling back to env", configPath)
	}
	return &config, nil
}

func createConfig(config *Config) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return errors.New("couldn't locate home directory")
	}
	var buffer bytes.Buffer
	if err := toml.NewEncoder(&buffer).Encode(config); err != nil {
		return err
	}
	if err := os.WriteFile(path.Join(home, CONFIG_DIR, CONFIG_FILENAME), buffer.Bytes(), 0500); err != nil {
		return err
	}
	return nil
}

func GetConfig() (*Config, error) {
	result, err := getConfigFileIfExists()
	if err != nil {
		return nil, err
	}

	return result, nil
}
