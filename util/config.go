package util

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/trevor-atlas/zilla/constants"
	"github.com/trevor-atlas/zilla/logger"
	"log"
	"os"
	"path"
)

type Jiraconf struct {
	Username     string `toml:"username,omitempty"`
	Apikey       string `toml:"apikey,omitempty"`
	Orgname      string `toml:"orgname,omitempty"`
	CustomDomain string `toml:"customDomain,omitempty"`
	AccessToken  string `toml:"accessToken,omitempty"`
}

type ConfigData struct {
	Jira  Jiraconf `toml:"jira,omitempty"`
	IsDev bool     `toml:"isDev,omitempty"`
}

type Zilla struct {
	Config *ConfigData
	Info   *log.Logger
	Err    *log.Logger
}

func New() *Zilla {
	app := new(Zilla)
	LogInfo, LogErr := logger.GetLoggers()
	app.Info = LogInfo
	app.Err = LogErr
	if config, err := app.GetConfig(); err != nil {
		app.Config = new(ConfigData)
	} else {
		app.Config = config
	}

	return app
}

func (a *Zilla) getConfigFileIfExists() (*ConfigData, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		a.Err.Println("error attempting to locate home directory")
		return nil, errors.New("couldn't locate home directory")
	}
	var conf ConfigData
	configPath := path.Join(home, constants.CONFIG_DIR, constants.CONFIG_FILENAME)
	_, err = toml.DecodeFile(configPath, &conf)
	if err != nil {
		a.Err.Printf("error parsing config file: %#v", err)
		return nil, fmt.Errorf("error parsing config at \"%v\" ", configPath)
	}
	a.Info.Printf("successfully loaded config at \"%s\"", configPath)
	return &conf, nil
}

func (a *Zilla) createConfig(config *ConfigData) error {
	home, err := os.UserHomeDir()
	if err != nil {
		a.Err.Println("couldn't locate home directory")
		return errors.New("couldn't locate home directory while attempting to create config file")
	}
	var buffer bytes.Buffer
	if err := toml.NewEncoder(&buffer).Encode(config); err != nil {
		a.Err.Println("error encoding config file while attempting to create it")
		return err
	}
	if err := os.WriteFile(path.Join(home, constants.CONFIG_DIR, constants.CONFIG_FILENAME), buffer.Bytes(), 0500); err != nil {
		a.Err.Println("error creating config file")
		return err
	}
	a.Info.Println("successfully created config file")
	return nil
}

func (a *Zilla) GetConfig() (*ConfigData, error) {
	result, err := a.getConfigFileIfExists()
	if err != nil {
		return nil, err
	}

	return result, nil
}
