package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
)

// WriteConfig creates the .uipo directory and then writes the config.json.
func (c *Config) WriteConfig() error {
	rawConfig, err := json.MarshalIndent(c.ConfigFile, "", "  ")
	if err != nil {
		return err
	}

	dir := configDirectory()
	err = os.MkdirAll(dir, 0700)
	if err != nil {
		return err
	}

	sig := make(chan os.Signal, 10)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, os.Interrupt)
	defer signal.Stop(sig)

	tempConfigFile, err := ioutil.TempFile(dir, "temp-config")
	if err != nil {
		return err
	}
	tempConfigFile.Close()
	tempConfigFileName := tempConfigFile.Name()

	go catchSignal(sig, tempConfigFileName)

	err = ioutil.WriteFile(tempConfigFileName, rawConfig, 0600)
	if err != nil {
		return err
	}

	return os.Rename(tempConfigFileName, ConfigFilePath())
}

// catchSignal tries to catch SIGHUP, SIGINT, SIGKILL, SIGQUIT and SIGTERM, and
// Interrupt for removing temporarily created config files before the program
// ends.
func catchSignal(sig chan os.Signal, tempConfigFileName string) {
	<-sig
	_ = os.Remove(tempConfigFileName)
	os.Exit(2)
}
