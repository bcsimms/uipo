package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/bcsimms/uipo/commands"
	"github.com/bcsimms/uipo/config"
	"github.com/bcsimms/uipo/util"
	"github.com/jessevdk/go-flags"
)

// CommandList is our list of supported commands
type CommandList struct {
	Verbose       bool                      `short:"v" long:"verbose" hidden:"true" description:"Run in verbose mode.  Outputs debug level interaction details"`
	Unsafe        bool                      `long:"unsafe" description:"Unsafe mode, disables endpoint certificate verification"`
	Authenticate  commands.CmdAuthenticate  `command:"authenticate" description:"Authenticate to UiPath Orchestrator"`
	PlatformSetup commands.CmdPlatformSetup `command:"platform-setup" description:"Used to setup UiPath Platform default values"`
	Robots        commands.CmdRobots        `command:"robots" description:"List Robots in current tenant"`
	Folders       commands.CmdGetFolders    `command:"folders" description:"List folders for current user"`
	UploadPackage commands.CmdUploadPackage `command:"push" description:"Upload a new package to Orchestrator"`
	AddQueueItem  commands.CmdAddQueueItem  `command:"addq" description:"Add an item to a queue"`
}

var cmds CommandList

func init() {
	log.SetFlags(log.Ldate | log.Lmicroseconds)
}

func main() {

	parser := flags.NewParser(&cmds, flags.Default)
	// Assign our custom execution handler
	parser.CommandHandler = executionWrapper
	// Parse command line and execute the provided command

	cmdLineArgs := os.Args[1:]

	parser.ParseArgs(cmdLineArgs)

	os.Exit(0)

}

func executionWrapper(cmd flags.Commander, args []string) error {

	uipoConfig, configErr := config.LoadConfig()
	if configErr != nil {
		return configErr
	}
	util.LogInfo("Configuration Loaded")

	defer func() {
		configWriteErr := uipoConfig.WriteConfig()
		if configWriteErr != nil {
			fmt.Fprintf(os.Stderr, "Error writing config: %s", configWriteErr.Error())
		}
		util.LogDebug("Completed deferred config file write operation")
	}()

	if cmds.Verbose {
		uipoConfig.GlobalFlgs.Verbose = true
		util.LogLevel = util.LogLevelDebug
	}

	// If running with --unsafe turn off certificate verification
	// Can be used if the endpoint includes a cert with mis-matched host name and assignee
	if cmds.Unsafe {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		uipoConfig.GlobalFlgs.Unsafe = true
		util.LogDebug("Running in Unsafe Mode")
	}

	if extendedCmd, ok := cmd.(commands.ExtendedCommander); ok {

		err := extendedCmd.Setup(uipoConfig)
		if err != nil {
			return err
		}

		return extendedCmd.Execute(args)

	}

	return errors.New("Command does not conform to ExtendedCommander")
}
