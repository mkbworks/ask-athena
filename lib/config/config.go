package config

import (
	"errors"
	"path/filepath"
	"runtime"
)

// Absolute path of the BIND file that contains the root DNS server details.
var RootServerFilePath string
// Absolute path of the BIND file that contains all the cached RRs.
var CacheFilePath string

//Sets up the inital configuration required to create a resolver instance.
func SetupConfig() error {
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
	return nil
}