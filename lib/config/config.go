package config

import (
	"os"
	"path/filepath"
	"time"
	"fmt"
)

var RootServerFilePath string
var CacheFilePath string
var LogFilePath string

func SetupConfig() error {
	CurrentDirectory, err := os.Getwd()
	if err != nil {
		return err
	}

	CacheFilePath = filepath.Join(CurrentDirectory, "assets", "resolver-cache.conf")
	RootServerFilePath = filepath.Join(CurrentDirectory, "assets", "root-servers.conf")
	CurrentTime := time.Now().UTC()
	LogFileName := fmt.Sprintf("%d%02d%02d%02d%02d%02d.log", CurrentTime.Year(), CurrentTime.Month(), CurrentTime.Day(), CurrentTime.Hour(), CurrentTime.Minute(), CurrentTime.Second())
	LogFilePath = filepath.Join(CurrentDirectory, LogFileName)
	return nil
}