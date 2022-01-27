package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"io/ioutil"
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
	CONFIG_FILENAME = "config.json"
	CACHE_FILENAME  = "cache.json"
)

var style = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FAFAFA")).
	Background(lipgloss.Color("#7D56F4")).
	Padding(1)

// TODO: should config files be created automatically? or should we fall back to finding values in the environment?
func getConfigFileIfExists() (*Config, error) {
	config := Config{}
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, errors.New("couldn't locate home directory, falling back to env based config")
	}
	xdgHomeExists := Exists(path.Join(home, CONFIG_DIR))
	if !xdgHomeExists {
		os.Create(path.Join(home, CONFIG_DIR))
		fmt.Println(style.Render(fmt.Sprintf("Created XDG home at ~/%v", CONFIG_DIR)))
	}

	configPath := path.Join(home, ".config/zilla/zilla.json")
	fmt.Println(configPath)
	jsonFile, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config at \"%v\" falling back to env", configPath)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		return nil, fmt.Errorf("error parsing config at \"%v\" falling back to env", configPath)
	}
	return &config, nil
}

func GetConfig() Config {
	result, err := getConfigFileIfExists()
	if err != nil {
		fmt.Println(err)
	}
	if result != nil {
		fmt.Println("using file config")
		return *result
	}

	fmt.Println("using env config")
	config := Config{}
	env := getenv()
	config.IsDev = env["ZILLA_IS_DEV"] == "true"
	config.Jira = Jiraconf{
		Username:     env["ZILLA_JIRA_USERNAME"],
		Apikey:       env["ZILLA_JIRA_APIKEY"],
		Orgname:      env["ZILLA_JIRA_ORG_NAME"],
		CustomDomain: env["ZILLA_JIRA_CUSTOM_DOMAIN"],
		AccessToken:  env["ZILLA_JIRA_ACCESS_TOKEN"],
	}
	return config
}
