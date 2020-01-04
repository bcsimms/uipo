// +build windows

package config

import (
	"os"
	"path/filepath"
)

// ConfigFilePath returns the location of the config file
func ConfigFilePath() string {
	return filepath.Join(configDirectory(), "config.json")
}

func configDirectory() string {
	return filepath.Join(homeDirectory(), ".uipo")
}

func homeDirectory() string {
	var homeDir string
	switch {
	case os.Getenv("UIPO_HOME") != "":
		homeDir = os.Getenv("UIPO_HOME")
	case os.Getenv("HOMEDRIVE")+os.Getenv("HOMEPATH") != "":
		homeDir = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
	default:
		homeDir = os.Getenv("USERPROFILE")
	}
	return homeDir
}
