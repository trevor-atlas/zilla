package logger

import (
	"fmt"
	"github.com/trevor-atlas/zilla/util"
	"log"
	"os"
	"path"
)

var CommonLog *log.Logger
var ErrorLog *log.Logger

func GetLoggers() (*log.Logger, *log.Logger) {
	home, err := os.UserHomeDir()
	openLogfile, err := os.OpenFile(path.Join(home, util.CONFIG_DIR, util.LOG_FILENAME), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error opening logfile:", err)
		os.Exit(1)
	}
	CommonLog = log.New(openLogfile, "Common Logger:\t", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLog = log.New(openLogfile, "Error Logger:\t", log.Ldate|log.Ltime|log.Lshortfile)
	return CommonLog, ErrorLog
}
