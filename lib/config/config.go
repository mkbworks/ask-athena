package config

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"time"
)

// Absolute path of the BIND file that contains the root DNS server details.
var RootServerFilePath string
// Absolute path of the BIND file that contains all the cached RRs.
var CacheFilePath string
// Absolute path of the log file that will be generated.
var LogFilePath string

//Sets up the inital configuration required to create a resolver instance.
func SetupConfig(LogFilesDirectory string) error {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return errors.New("unable to fetch current file path")
	}

	completeFilePath, err := filepath.Abs(file)
	if err != nil {
		return err
	}

	CurrentDirectory := filepath.Dir(completeFilePath)
	CacheFilePath = filepath.Join(CurrentDirectory, "resolver-cache.conf")
	RootServerFilePath = filepath.Join(CurrentDirectory, "root-servers.conf")

	if !filepath.IsAbs(LogFilesDirectory) {
		return errors.New("log file directory path is not an absolute path")
	}

	CurrentTime := time.Now().UTC()
	LogFileName := fmt.Sprintf("DNS_LOGS_%d%02d%02d_%02d%02d%02d.log", CurrentTime.Year(), CurrentTime.Month(), CurrentTime.Day(), CurrentTime.Hour(), CurrentTime.Minute(), CurrentTime.Second())
	LogFilePath = filepath.Join(LogFilesDirectory, LogFileName)
	return nil
}