package commands

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/bcsimms/uipo/config"
	"github.com/bcsimms/uipo/util"
)

const upEndPointURI string = "/odata/Processes/UiPath.Server.Configuration.OData.UploadPackage"

// CmdUploadPackage represents the flags this command support
type CmdUploadPackage struct {
	APIEndpoint        string `short:"e" long:"api-endpoint" description:"API endpoint (e.g. https://api.example.com)"`
	AccountLogicalName string `short:"a" long:"alname" description:"Account Logical Name - Used for UiPath Platform Installations"`
	ServiceLogicalName string `short:"s" long:"slname" description:"Service Logical Name - Used for UiPath Platform Installations"`
	PackagePath        string `short:"p" long:"package" description:"Full path to NuGet package file"`

	Config Config
}

type uploadRespWrapper struct {
	ODataContext string       `json:"@odata.context"`
	UploadResp   []uploadResp `json:"value"`
}

type uploadResp struct {
	Key    string `json:"Key"`
	Status string `json:"Status"`
	Body   string `json:"Body"`
}

// Setup is the standard setup function
func (cmd *CmdUploadPackage) Setup(conf Config) error {

	cmd.Config = conf

	// If we have flags passed in with the command, use those and store them as
	//   the new cached values in our config
	// If both the flags and cached values are empty string, the validateFlags func
	//   will catch that scenario
	if cmd.APIEndpoint == "" {
		cmd.APIEndpoint = cmd.Config.GetAPIEndpoint()
	} else {
		cmd.Config.SetAPIEndpoint(cmd.APIEndpoint)
	}
	if cmd.AccountLogicalName == "" {
		cmd.AccountLogicalName = cmd.Config.GetAccountLogicalName()
	} else {
		cmd.Config.SetAccountLogicalName(cmd.AccountLogicalName)
	}
	if cmd.ServiceLogicalName == "" {
		cmd.ServiceLogicalName = cmd.Config.GetServiceLogicalName()
	} else {
		cmd.Config.SetServiceLogicalName(cmd.ServiceLogicalName)
	}

	return nil

}

// Execute is the main entry point for this command
func (cmd *CmdUploadPackage) Execute(args []string) error {

	// Validate incoming flags
	err := cmd.validateFlags()

	if err != nil {
		return err
	}

	// Open our package file
	file, err := os.Open(cmd.PackagePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Setup our multipart writer and create the request body using the file supplied
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(cmd.PackagePath))
	if err != nil {
		return err
	}
	fileSize, err := io.Copy(part, file)
	err = writer.Close()
	if err != nil {
		return err
	}

	var endpoint string

	if cmd.Config.GetEndpointType() == config.EndpointTypeHosted {
		endpoint = cmd.APIEndpoint + "/" + cmd.Config.GetAccountLogicalName() + "/" + cmd.Config.GetServiceLogicalName() + upEndPointURI
	} else if cmd.Config.GetEndpointType() == config.EndpointTypeOnPremise {
		endpoint = cmd.APIEndpoint + upEndPointURI
	} else {
		return errors.New("Invalid Endpoint Type in cached config.  Reauthenticate to reset")
	}

	client := http.Client{}
	req, err := http.NewRequest("POST", endpoint, body)
	if err != nil {
		return err
	}

	// Add our required request headers
	req.Header.Add("X-UIPATH-TenantName", cmd.Config.GetServiceLogicalName())
	req.Header.Add("Authorization", "Bearer "+cmd.Config.GetAccessToken())
	req.Header.Add("Content-Type", writer.FormDataContentType())

	// Use the HTTPHelper to make our API call
	respBody, err := util.HTTPHelper(client, req)
	if err != nil {
		return err
	}

	apiResp := uploadRespWrapper{}
	jsonErr := json.Unmarshal(respBody, &apiResp)
	if jsonErr != nil {
		return jsonErr
	}

	fmt.Println("Package uploaded successfully")
	fmt.Println("")
	fmt.Println("  Package File: ", apiResp.UploadResp[0].Key)
	fmt.Println("     File Size: ", fileSize)
	fmt.Println("        Status: ", apiResp.UploadResp[0].Status)
	fmt.Println("")

	return nil

}

func (cmd *CmdUploadPackage) validateFlags() error {

	if cmd.PackagePath == "" {
		return errors.New("A package file is required")
	}
	if cmd.APIEndpoint == "" {
		return errors.New("An API end point is required and a value was not found in the cached config")
	}

	return nil

}
