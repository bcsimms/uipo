package commands

// CmdPlatformSetup represents the flags this command supports
type CmdPlatformSetup struct {
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

	return nil
}

//If a new foldername is provided, call the UiPath API for populate the Folder struct
// FolderID is required for most API calls
func resolveFolderID() error {

	return nil
}
