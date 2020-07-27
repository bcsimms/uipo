package config

import "strconv"

//Config is the main configuration object used throughout the UIPO CLI
type Config struct {
	// ConfigFile stores the configuration from the .uipo/config
	ConfigFile JSONConfig

	// ENV stores the configuration from os.ENV
	ENV EnvOverride

	// detectedSettings are settings detected when the config is loaded
	detectedSettings detectedSettings

	GlobalFlgs globalFlgs
}

// JSONConfig is the representation of the contents of our configuration file
// It stores temporary authentication details in addition to other info that might persist
// between API calls
type JSONConfig struct {
	ConfigVersion         int    `json:"ConfigVersion"`
	AccessToken           string `json:"AccessToken"`
	APIVersion            string `json:"APIVersion"`
	AsyncTimeout          int    `json:"AsyncTimeout"`
	AuthorizationEndpoint string `json:"AuthorizationEndpoint"`
	APIEndpoint           string `json:"APIEndpoint"`
	EndpointType          string `json:"EndpointType"`
	TargetedTenant        Tenant `json:"TenantFields"`
	TargetedFolder        Folder `json:"FolderFields"`
	RefreshToken          string `json:"RefreshToken"`
	Trace                 string `json:"Trace"`
	AccountLogicalName    string `json:"AccountLogicalName"`
	ServiceLogicalName    string `json:"ServiceLogicalName"`
	ClientID              string `json:"ClientID"`
}

// Tenant is the representation of a Tenant object in UiPath Orchestrator
//  This will store the default tenant used for API calls
type Tenant struct {
	Name string `json:"Name"`
	ID   string `json:"ID"`
	Key  string `json:"Key"`
}

// Folder is the representation of a Folder in UiPath Orchestrator
type Folder struct {
	DisplayName        string `json:"DisplayName"`
	FullyQualifiedName string `json:"FullyQualifiedName":`
	Description        string `json:"Description"`
	ParentID           int    `json:"ParentId"`
	ID                 int    `json:"Id"`
}

//EnvOverride is the represetnation of environment variables that might override
//  other config
type EnvOverride struct {
	BinaryName   string
	UIPOHome     string
	UIPOPassword string
	UIPOUsername string
	HTTPSProxy   string
}

const EndpointTypeHosted string = "Hosted"
const EndpointTypeOnPremise string = "OnPremise"

type detectedSettings struct {
	currentDirectory string
}

type globalFlgs struct {
	Unsafe  bool
	Verbose bool
}

func (config *Config) GetConfigVersion() string {
	return strconv.Itoa(config.ConfigFile.ConfigVersion)
}
func (config *Config) GetAccessToken() string {
	return config.ConfigFile.AccessToken
}

func (config *Config) GetAPIVersion() string {
	return config.ConfigFile.APIVersion
}

func (config *Config) GetAPIEndpoint() string {
	return config.ConfigFile.APIEndpoint
}

// SetAccessToken sets the current access token.
func (config *Config) SetAccessToken(accessToken string) {
	config.ConfigFile.AccessToken = accessToken
}

func (config *Config) GetBinaryVersion() string {
	panic("not implemented")
}

func (config *Config) GetRefreshToken() string {
	return config.ConfigFile.RefreshToken
}

func (config *Config) SetRefreshToken(token string) {
	config.ConfigFile.RefreshToken = token
}

//SetAPIVersion - comes from API call to Orchestrator
func (config *Config) SetAPIVersion(apiVerison string) {
	config.ConfigFile.APIVersion = apiVerison
}

// SetAuthorizationEndpoint sets the authoritcation endpoint for future use
func (config *Config) SetAuthorizationEndpoint(e string) {
	config.ConfigFile.AuthorizationEndpoint = e
}

func (config *Config) GetAuthorizationEndpoint() string {
	return config.ConfigFile.AuthorizationEndpoint
}

//SetAPIEndpoint sets the API endpoint for future use
func (config *Config) SetAPIEndpoint(apiEndpoint string) {
	config.ConfigFile.APIEndpoint = apiEndpoint
}

func (config *Config) GetEndpointType() string {
	return config.ConfigFile.EndpointType
}

func (config *Config) SetEndpointType(endpointType string) {
	config.ConfigFile.EndpointType = endpointType
}

// BinaryName returns the running name of the UIPO CLI
func (config *Config) GetBinaryName() string {
	return config.ENV.BinaryName
}

// UIPOPassword returns the value of the "UIPO_PASSWORD" environment variable.
func (config *Config) GetUIPOPassword() string {
	return config.ENV.UIPOPassword
}

// UIPOUsername returns the value of the "UIPO_USERNAME" environment variable.
func (config *Config) GetUIPOUsername() string {
	return config.ENV.UIPOUsername
}

func (config *Config) GetAccountLogicalName() string {
	return config.ConfigFile.AccountLogicalName
}

func (config *Config) SetAccountLogicalName(name string) {
	config.ConfigFile.AccountLogicalName = name
}

func (config *Config) GetServiceLogicalName() string {
	return config.ConfigFile.ServiceLogicalName
}

func (config *Config) SetServiceLogicalName(name string) {
	config.ConfigFile.ServiceLogicalName = name
}

func (config *Config) GetClientID() string {
	return config.ConfigFile.ClientID
}

func (config *Config) SetClientID(id string) {
	config.ConfigFile.ClientID = id
}

func (config *Config) GetFolderID() int {
	return config.ConfigFile.TargetedFolder.ID
}

func (config *Config) SetFolderID(Id int) {
	config.ConfigFile.TargetedFolder.ID = Id
}

func (config *Config) GetFolderFQN() string {
	return config.ConfigFile.TargetedFolder.FullyQualifiedName
}

func (config *Config) SetFolderFQN(fqn string) {
	config.ConfigFile.TargetedFolder.FullyQualifiedName = fqn
}

func (config *Config) GetFolderName() string {
	return config.ConfigFile.TargetedFolder.DisplayName
}

func (config *Config) SetFolderName(name string) {
	config.ConfigFile.TargetedFolder.DisplayName = name
}

func (config *Config) GetFolderDescription() string {
	return config.ConfigFile.TargetedFolder.Description
}

func (config *Config) SetFolderDescription(desc string) {
	config.ConfigFile.TargetedFolder.Description = desc
}

func (config *Config) GetFolderParentID() int {
	return config.ConfigFile.TargetedFolder.ParentID
}

func (config *Config) SetFolderParentID(id int) {
	config.ConfigFile.TargetedFolder.ParentID = id
}

func (config *Config) IsUnsafeMode() bool {
	return config.GlobalFlgs.Unsafe
}

func (config *Config) IsVerboseMode() bool {
	return config.GlobalFlgs.Verbose
}
