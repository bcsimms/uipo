package commands

import (
	"fmt"
	"strconv"
)

// CmdPlatformSetup represents the flags this command supports
type CmdPlatformSetup struct {
	View                 bool   `short:"v" long:"view" description:"Displays the current setup (contents of the .uipo config file)"`
	APIEndpoint          string `short:"e" long:"api-endpoint" description:"API endpoint (e.g. https://api.example.com)"`
	AuthticationEndpoint string `short:"u" long:"auth-endpoint" description:"Endpoint used to generate a new bearer token"`
	RefreshToken         string `short:"r" long:"refreh-token" description:"Refresh Token - Used for UiPath Platform Installations"`
	AccountLogicalName   string `short:"a" long:"alname" description:"Account Logical Name - Used for UiPath Platform Installations"`
	ServiceLogicalName   string `short:"s" long:"slname" description:"Service Logical Name - Used for UiPath Platform Installations"`
	ClientID             string `short:"c" long:"client-id" default:"5v7PmPJL6FOGu6RB8I1Y4adLBhIwovQN" descritpion:"Client ID - Used for UiPath Platform Installations.  Should not need to be overridden"`
	FolderName           string `short:"f" long:"folder" description:"The UiPath folder to be used for subsequent API operations"`

	Config Config
}

// Setup is the standard setup function
func (cmd *CmdPlatformSetup) Setup(config Config) error {

	cmd.Config = config

	return nil
}

// Execute is the main entry point for this command
func (cmd *CmdPlatformSetup) Execute(args []string) error {

	if cmd.View {
		fmt.Println("")
		fmt.Println("Current UIPO Configuration")
		fmt.Println("      Config Version: " + cmd.Config.GetConfigVersion())
		fmt.Println("         API Version: " + cmd.Config.GetAPIVersion())
		fmt.Println("            API Type: " + cmd.Config.GetEndpointType())
		fmt.Println("       Auth Endpoint: " + cmd.Config.GetAuthorizationEndpoint())
		fmt.Println("         API Enpoint: " + cmd.Config.GetAPIEndpoint())
		fmt.Println("      Default Folder: " + cmd.Config.GetFolderFQN() + "; ID: " + strconv.Itoa(cmd.Config.GetFolderID()))
		fmt.Println("Account Logical Name: " + cmd.Config.GetAccountLogicalName())
		fmt.Println("Service Logical Name: " + cmd.Config.GetServiceLogicalName())
		fmt.Println("          User Token: " + cmd.Config.GetRefreshToken())
		fmt.Println("           Client ID: " + cmd.Config.GetClientID())
		fmt.Println("")

	} else {
		if cmd.APIEndpoint != "" {
			cmd.Config.SetAPIEndpoint(cmd.APIEndpoint)
		}
		if cmd.AuthticationEndpoint != "" {
			cmd.Config.SetAuthorizationEndpoint(cmd.AuthticationEndpoint)
		}
		if cmd.RefreshToken != "" {
			cmd.Config.SetRefreshToken(cmd.RefreshToken)
		}
		if cmd.AccountLogicalName != "" {
			cmd.Config.SetAccountLogicalName(cmd.AccountLogicalName)
		}
		if cmd.ServiceLogicalName != "" {
			cmd.Config.SetServiceLogicalName(cmd.ServiceLogicalName)
		}
		if cmd.ClientID != "" {
			cmd.Config.SetClientID(cmd.ClientID)
		}
		if cmd.FolderName != "" {
			cmd.Config.SetFolderName(cmd.FolderName)
		}

	}

	return nil
}

//If a new foldername is provided, call the UiPath API to populate the Folder struct
// FolderID is required for most API calls
func resolveFolderID() error {

	return nil
}
