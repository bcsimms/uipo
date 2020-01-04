package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/ekaju-of-co/uipo/config"
)

// CmdRobots represents the flags supported by this command
type CmdRobots struct {
	APIEndpoint        string `short:"e" long:"endpoint" description:"API enpoint base URL"`
	AccountLogicalName string `short:"a" long:"alname" description:"Account Logical Name - Used for UiPath Platform Installations"`
	ServiceLogicalName string `short:"s" long:"slname" description:"Service Logical Name - Used for UiPath Platform Installations"`

	Config Config
}

type odataResp struct {
	ODataContext string   `json:"@odata.context"`
	ODataCount   int      `json:"@odata.count"`
	Robots       []robots `json:"value"`
}
type robots struct {
	LicenseKey  string `json:"LicenseKey"`
	MachineName string `json:"MachineName"`
	MachineID   int    `json:"MachineId"`
	Name        string `json:"Name"`
	Version     string `json:"Version"`
}

// Setup is the standard setup function
func (cmd *CmdRobots) Setup(config Config) error {

	cmd.Config = config

	if cmd.Config.GetAPIEndpoint() != "" {
		cmd.APIEndpoint = cmd.Config.GetAPIEndpoint()
	}
	return nil

}

// Execute is the main entry point for this command
func (cmd *CmdRobots) Execute(args []string) error {

	err := cmd.validateFlags()

	if err != nil {
		return err
	}

	var robotsEndpoint string
	if cmd.Config.GetEndpointType() == config.EndpointTypeHosted {
		robotsEndpoint = cmd.APIEndpoint + "/" + cmd.Config.GetAccountLogicalName() + "/" + cmd.Config.GetServiceLogicalName() + "/odata/Robots"
	} else if cmd.Config.GetEndpointType() == config.EndpointTypeOnPremise {
		robotsEndpoint = cmd.APIEndpoint + "/odata/Robots"
	} else {
		errors.New("Invalid Endpoint Type in cached config.  Reauthenticate to reset")
	}

	client := http.Client{}
	req, err := http.NewRequest("GET", robotsEndpoint, nil)
	req.Header.Add("X-UIPATH-TenantName", cmd.Config.GetServiceLogicalName())
	req.Header.Add("X-UIPATH-OrganizationUnitId", strconv.Itoa(cmd.Config.GetFolderID()))
	req.Header.Add("Authorization", "Bearer "+cmd.Config.GetAccessToken())

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error with Client Request")
		return err
	}

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return readErr
	}
	defer resp.Body.Close()

	//fmt.Println(string(body))

	odataResp := odataResp{}

	jsonErr := json.Unmarshal(body, &odataResp)
	if jsonErr != nil {
		return jsonErr
	}

	if odataResp.ODataCount == 0 {
		fmt.Println("No Robots returned")
		fmt.Println("")
	} else {
		for _, element := range odataResp.Robots {
			fmt.Println("       Robot Name: ", element.Name)
			fmt.Println("       Machine ID: ", element.MachineID)
			fmt.Println("     Machine Name: ", element.MachineName)
			fmt.Println("  Machine Version: ", element.Version)
			fmt.Println("      License Key: ", element.LicenseKey)
			fmt.Println("")
		}

	}

	//fmt.Println(string(body))
	//fmt.Println(odataResp.Robots[0].Name)

	return nil
}

func (cmd *CmdRobots) validateFlags() error {

	return nil

}
