package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

// LoadConfig loads the config from the .cf/config.json and os.ENV. If the
// config.json does not exists, it will use a default config in it's place.
// Takes in an optional FlagOverride, will only use the first one passed, that
// can override the given flag values.
//
// The '.cf' directory will be read in one of the following locations on UNIX
// Systems:
//   1. $CF_HOME/.cf if $CF_HOME is set
//   2. $HOME/.cf as the default
//
// The '.cf' directory will be read in one of the following locations on
// Windows Systems:
//   1. CF_HOME\.cf if CF_HOME is set
//   2. HOMEDRIVE\HOMEPATH\.cf if HOMEDRIVE or HOMEPATH is set
//   3. USERPROFILE\.cf as the default
func LoadConfig() (*Config, error) {
	err := removeOldTempConfigFiles()
	if err != nil {
		return nil, err
	}

	configFilePath := ConfigFilePath()

	config := Config{
		ConfigFile: JSONConfig{
			ConfigVersion: 1,
			ClientID:      DefaultClientID,
		},
	}

	var jsonError error

	if _, err = os.Stat(configFilePath); err == nil || !os.IsNotExist(err) {
		var file []byte
		file, err = ioutil.ReadFile(configFilePath)
		if err != nil {
			return nil, err
		}

		if len(file) != 0 {

			var configFile JSONConfig
			err = json.Unmarshal(file, &configFile)
			if err != nil {
				return nil, err
			}
			config.ConfigFile = configFile
		}
	}

	config.ENV = EnvOverride{
		BinaryName:   filepath.Base(os.Args[0]),
		UIPOPassword: os.Getenv("UIPO_PASSWORD"),
		UIPOUsername: os.Getenv("UIPO_USERNAME"),
		UIPOHome:     os.Getenv("UIPO_HOME"),
		HTTPSProxy:   os.Getenv("https_proxy"),
	}

	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	config.detectedSettings = detectedSettings{
		currentDirectory: pwd,
	}

	config.GlobalFlgs = globalFlgs{
		Unsafe:  false,
		Verbose: false,
	}

	return &config, jsonError
}

func removeOldTempConfigFiles() error {
	oldTempFileNames, err := filepath.Glob(filepath.Join(configDirectory(), "temp-config?*"))
	if err != nil {
		return err
	}

	for _, oldTempFileName := range oldTempFileNames {
		err = os.Remove(oldTempFileName)
		if err != nil {
			return err
		}
	}

	return nil
}
