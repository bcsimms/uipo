package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/ekaju-of-co/uipo/config"
	"github.com/ekaju-of-co/uipo/util"
)

// CmdGetFolders represetns the flags this command supports
type CmdGetFolders struct {
	APIEndpoint        string `short:"e" long:"api-endpoint" description:"API endpoint (e.g. https://api.example.com)"`
	AccountLogicalName string `short:"a" long:"alname" description:"Account Logical Name - Used for UiPath Platform Installations"`
	ServiceLogicalName string `short:"s" long:"slname" description:"Service Logical Name - Used for UiPath Platform Installations"`
	FilteredFolders    string `short:"f" long:"filter" description:"Run operation as the filtered version using the String as the filter criterial"`
	Skip               string `short:"x" long:"skip" default:"0" description:"For use when using the filtered version.  Will skip some number of entries"`
	Take               string `short:"y" long:"take" default:"10" description:"For use when using the filtered version.  Determines how many results to return"`
	SetDefaultFolder   bool   `short:"d" long:"set-default" description:"Setting this flag to true will persist the fisrt folder result as the default"`

	Config Config
}

type filteredFoldersResp struct {
	FolderEntries []filteredFolderEntry `json:"PageItems"`
	Count         int                   `json:"Count"`
}

type filteredFolderEntry struct {
	DisplayName        string `json:"DisplayName"`
	FullyQualifiedName string `json:"FullyQualifiedName"`
	Description        string `json:"Description"`
	ProvisionType      string `json:"ProvisionType"`
	PermissionModel    string `json:"PermissionModel"`
	ParentID           int    `json:"ParentId"`
	ID                 int    `json:"Id"`
}

type allFolderEntry struct {
	IsSelectable       bool   `json:"IsSelectable"`
	HasChildren        bool   `json:"HasChildren"`
	Level              int    `json:"Level"`
	DisplayName        string `json:"DisplayName"`
	FullyQualifiedName string `json:"FullyQualifiedName"`
	Description        string `json:"Description"`
	ProvisionType      string `json:"ProvisionType"`
	PermissionModel    string `json:"PermissionModel"`
	ParentID           int    `json:"ParentId"`
	ID                 int    `json:"Id"`
}

// Setup is the standard setup function - override of the go-flags interface function
func (cmd *CmdGetFolders) Setup(config Config) error {

	cmd.Config = config

	// If we have input flag overrides, use those instead of cached values
	if cmd.APIEndpoint == "" {
		cmd.APIEndpoint = cmd.Config.GetAPIEndpoint()
	}
	if cmd.AccountLogicalName == "" {
		cmd.AccountLogicalName = cmd.Config.GetAccountLogicalName()
	}
	if cmd.ServiceLogicalName == "" {
		cmd.ServiceLogicalName = cmd.Config.GetServiceLogicalName()
	}

	return nil
}

// Execute is the main entry poing for the command - override of the go-flags interface function
func (cmd *CmdGetFolders) Execute(args []string) error {

	var err error
	if cmd.FilteredFolders != "" {
		err = cmd.getFilteredFolders()
	} else {
		err = cmd.getAllFolders()
	}

	// If we have a successful execution, do some housekeeping
	if err == nil {
		// Save arguments to config
		if cmd.APIEndpoint != "" {
			cmd.Config.SetAPIEndpoint(cmd.APIEndpoint)
		}
	}
	return err
}

// Usage is the override for the go-flags interface function
func (cmd *CmdGetFolders) Usage() string {
	usageStr := "[-e API Endpoint URL] [-f Filter Condition] [-x Skip entries] [-y How many entries to return] [-d Sets first returned entry, if any, as default]"
	return usageStr
}

func (cmd *CmdGetFolders) getFilteredFolders() error {

	var foldersEndpoint string

	if cmd.Config.GetEndpointType() == config.EndpointTypeHosted {
		foldersEndpoint = cmd.APIEndpoint + "/" + cmd.AccountLogicalName + "/" + cmd.ServiceLogicalName + "/api/FoldersNavigation/GetFoldersForCurrentUser"
	} else if cmd.Config.GetEndpointType() == config.EndpointTypeOnPremise {
		foldersEndpoint = cmd.APIEndpoint + "/api/FoldersNavigation/GetFoldersForCurrentUser"
	} else {
		return errors.New("Invalid Endpoint Type in cached config.  Reauthenticate to reset")
	}

	client := http.Client{}
	req, _ := http.NewRequest("GET", foldersEndpoint, nil)

	// Add our required request headers
	req.Header.Add("X-UIPATH-TenantName", cmd.ServiceLogicalName)
	req.Header.Add("Authorization", "Bearer "+cmd.Config.GetAccessToken())

	// Add our query string bits and pieces
	query := req.URL.Query()
	query.Add("searchText", cmd.FilteredFolders)
	query.Add("skip", cmd.Skip)
	query.Add("take", cmd.Take)
	req.URL.RawQuery = query.Encode()

	// Use the HTTPHelper to make our API call
	body, err := util.HTTPHelper(client, req)

	if err != nil {
		return err
	}

	// Process oour response
	apiResp := filteredFoldersResp{}
	jsonErr := json.Unmarshal(body, &apiResp)
	if jsonErr != nil {
		return jsonErr
	}

	// If we get a "-d" save the folder details we found
	if cmd.SetDefaultFolder {
		cmd.Config.SetFolderID(apiResp.FolderEntries[0].ID)
		cmd.Config.SetFolderFQN(apiResp.FolderEntries[0].FullyQualifiedName)
		cmd.Config.SetFolderDescription(apiResp.FolderEntries[0].Description)
		cmd.Config.SetFolderName(apiResp.FolderEntries[0].DisplayName)
		cmd.Config.SetFolderParentID(apiResp.FolderEntries[0].ParentID)
	}

	//fmt.Println(string(body))
	for _, element := range apiResp.FolderEntries {
		fmt.Println("           Folder ID: ", element.ID)
		fmt.Println("        Display Name: ", element.DisplayName)
		fmt.Println("Fully Qualified Name: ", element.FullyQualifiedName)
		fmt.Println("         Description: ", element.Description)
		fmt.Println("           Parent ID: ", element.ParentID)
		fmt.Println("")
	}

	return nil
}

func (cmd *CmdGetFolders) getAllFolders() error {

	var foldersEndpoint string

	if cmd.Config.GetEndpointType() == config.EndpointTypeHosted {
		foldersEndpoint = cmd.APIEndpoint + "/" + cmd.AccountLogicalName + "/" + cmd.ServiceLogicalName + "/api/FoldersNavigation/GetAllFoldersForCurrentUser"
	} else if cmd.Config.GetEndpointType() == config.EndpointTypeOnPremise {
		foldersEndpoint = cmd.APIEndpoint + "/api/FoldersNavigation/GetAllFoldersForCurrentUser"
	} else {
		return errors.New("Invalid Endpoint Type in cached config.  Reauthenticate to reset")
	}

	client := http.Client{}
	req, _ := http.NewRequest("GET", foldersEndpoint, nil)

	// Add our required request headers
	req.Header.Add("X-UIPATH-TenantName", cmd.ServiceLogicalName)
	req.Header.Add("Authorization", "Bearer "+cmd.Config.GetAccessToken())

	// Use the HTTPHelper to make our API call
	body, err := util.HTTPHelper(client, req)

	if err != nil {
		return err
	}

	// Process our response
	var apiResp []allFolderEntry
	jsonErr := json.Unmarshal(body, &apiResp)
	if jsonErr != nil {
		return jsonErr
	}

	// If we get a "-d" save the folder details we found
	if cmd.SetDefaultFolder {
		cmd.Config.SetFolderID(apiResp[0].ID)
		cmd.Config.SetFolderFQN(apiResp[0].FullyQualifiedName)
		cmd.Config.SetFolderDescription(apiResp[0].Description)
		cmd.Config.SetFolderName(apiResp[0].DisplayName)
		cmd.Config.SetFolderParentID(apiResp[0].ParentID)
	}

	//fmt.Println(string(body))
	for _, element := range apiResp {
		fmt.Println("           Folder ID: ", element.ID)
		fmt.Println("        Display Name: ", element.DisplayName)
		fmt.Println("Fully Qualified Name: ", element.FullyQualifiedName)
		fmt.Println("         Description: ", element.Description)
		fmt.Println("        Has Children: ", element.HasChildren)
		fmt.Println("           Parent ID: ", element.ParentID)
		fmt.Println("")
	}

	return nil
}
