package cache

import (
	"fmt"
	"github.com/trevor-atlas/zilla/jira"
	"github.com/trevor-atlas/zilla/util"
	"os"
	"path"
)

func GetCachedIssues() *jira.JiraIssues {
	// get cached issues from fs
	// I want to store these cached issues in a json file and then quickly read them to make the UI responsive faster
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("error reading filesystem")
	}
	issuesPath := path.Join(home, util.CONFIG_DIR, util.CACHE_FILENAME)
	os.Open(issuesPath)

}
