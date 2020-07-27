package commands

import (
	"github.com/jessevdk/go-flags"
)

//Config is the interface that wraps the Config Struct and is used
type Config interface {
	GetConfigVersion() string
	GetAccessToken() string
	GetAPIVersion() string
	GetAPIEndpoint() string
	GetAuthorizationEndpoint() string
	GetBinaryName() string
	GetBinaryVersion() string
	GetRefreshToken() string
	SetRefreshToken(token string)
	WriteConfig() error
	SetAPIVersion(string)
	SetAuthorizationEndpoint(string)
	SetAPIEndpoint(string)
	GetEndpointType() string
	SetEndpointType(string)
	SetAccessToken(string)
	GetUIPOPassword() string
	GetUIPOUsername() string
	GetAccountLogicalName() string
	SetAccountLogicalName(string)
	GetServiceLogicalName() string
	SetServiceLogicalName(string)
	GetClientID() string
	SetClientID(string)
	GetFolderID() int
	SetFolderID(int)
	GetFolderName() string
	SetFolderName(string)
	GetFolderFQN() string
	SetFolderFQN(string)
	GetFolderDescription() string
	SetFolderDescription(string)
	GetFolderParentID() int
	SetFolderParentID(int)
}

// ExtendedCommander is a type used for add a Setup function to commands
//  extending the go-flags command type
//  Right now setup only takes a config struct as an argument to bind the config to the command
type ExtendedCommander interface {
	flags.Commander
	Setup(Config) error
}
