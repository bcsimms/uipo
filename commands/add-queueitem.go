package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/bcsimms/uipo/config"
)

const aqiEndPointURI string = "/odata/Queues/UiPathODataSvc.AddQueueItem"

// CmdAddQueueItem represents the flags this command supports
type CmdAddQueueItem struct {
	APIEndpoint        string `short:"e" long:"api-endpoint" description:"API endpoint (e.g. https://api.example.com)"`
	AccountLogicalName string `short:"a" long:"alname" description:"Account Logical Name - Used for UiPath Platform Installations"`
	ServiceLogicalName string `short:"s" long:"slname" description:"Service Logical Name - Used for UiPath Platform Installations"`
	QueueName          string `short:"q" long:"queue" required:"true" description:"The name of the queue to which we will add a item"`
	Priority           string `short:"p" long:"priority" default:"Normal" choice:"Low" choice:"Normal" choice:"High" description:"The priority of the queue item."`
	Reference          string `short:"r" long:"reference" required:"true" description:"The reference identfier to assign to the queue item. Note: Some queues are configured to use unique reference ids"`
	SpecificContent    string `short:"c" long:"content" description:"Additional data for the queue item.  Must be provided as a single-quoted JSON string. (E.g. '{\"Attribue\":\"Value\"}')"`
	DueDate            string `short:"d" long:"deadline" description:"The UTC date before which the queue item should be processed. Must be provided as a JSON datetime in the format YYYY-MM-DDTHH:MM:SS.sssZ."`
	Postpone           string `long:"postpone" description:"The UTC date after which the queue item may be processed. Must be provided as a JSON datetime in the format YYYY-MM-DDTHH:MM:SS.sssZ."`

	Config Config
}

type requestBodyWrapper struct {
	ItemData requestBody `json:"itemData"`
}

type requestBody struct {
	Name            string `json:"Name"`
	Priority        string `json:"Priority"`
	SpecificContent string `json:"SpecificContent"`
	DeferDate       string `json:"DeferDate,omitempty"`
	DueDate         string `json:"DueDate,omitempty"`
	Reference       string `json:"Reference"`
}

type responseBody struct {
	ODataContext      string `json:"@odata.context"`
	QueueDefinitionID int    `json:"QueueDefinitionId"`
	Status            string `json:"Status"`
	Reference         string `json:"Reference"`
	CreationTime      string `json:"CreationTime"`
	ID                int    `json:"Id"`
}

type addQIErrorResponse struct {
	Message   string `json:"message"`
	ErrorCode int    `json:"errorCode"`
}

// Setup is called during execution to configure our command processing
func (cmd *CmdAddQueueItem) Setup(conf Config) error {

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
func (cmd *CmdAddQueueItem) Execute(args []string) error {

	// Validate incoming flags
	//err := cmd.validateFlags()

	//if err != nil {
	//	return err
	//}

	// Build our request body
	reqBody := requestBody{}
	reqBody.Name = cmd.QueueName
	reqBody.Reference = cmd.Reference
	reqBody.Priority = cmd.Priority
	reqBody.DueDate = cmd.DueDate
	reqBody.DeferDate = cmd.Postpone

	reqBodyWrapper := requestBodyWrapper{}
	reqBodyWrapper.ItemData = reqBody

	requestBody, err := json.Marshal(&reqBodyWrapper)
	if err != nil {
		return err
	}

	//Add the JSON specific content to the Request Body string
	requestBody = []byte(strings.Replace(string(requestBody), "\"SpecificContent\":\"\"", "\"SpecificContent\":"+cmd.SpecificContent, 1))

	var endpoint string

	if cmd.Config.GetEndpointType() == config.EndpointTypeHosted {
		endpoint = cmd.APIEndpoint + "/" + cmd.Config.GetAccountLogicalName() + "/" + cmd.Config.GetServiceLogicalName() + aqiEndPointURI
	} else if cmd.Config.GetEndpointType() == config.EndpointTypeOnPremise {
		endpoint = cmd.APIEndpoint + aqiEndPointURI
	} else {
		return errors.New("Invalid Endpoint Type in cached config.  Reauthenticate to reset")
	}

	client := http.Client{}
	req, err := http.NewRequest("POST", endpoint, strings.NewReader(string(requestBody)))
	if err != nil {
		return err
	}

	// Add our required request headers
	req.Header.Add("X-UIPATH-TenantName", cmd.Config.GetServiceLogicalName())
	req.Header.Add("X-UIPATH-OrganizationUnitId", strconv.Itoa(cmd.Config.GetFolderID()))
	req.Header.Add("Authorization", "Bearer "+cmd.Config.GetAccessToken())
	req.Header.Add("Content-Type", "application/json")

	respBody, err := client.Do(req)
	if err != nil {
		return err
	}

	body, readErr := ioutil.ReadAll(respBody.Body)
	if readErr != nil {
		return readErr
	}
	defer respBody.Body.Close()

	if respBody.StatusCode < 300 {
		apiResp := responseBody{}
		jsonErr := json.Unmarshal(body, &apiResp)
		if jsonErr != nil {
			return jsonErr
		}

		cmd.Config.SetAPIVersion(respBody.Header.Get("Api-Supported-Versions"))

		fmt.Println("Queue item created successfully")
		fmt.Println("")
		fmt.Println("      Item ID: ", apiResp.ID)
		fmt.Println("Creation Time: ", apiResp.CreationTime)
		fmt.Println("       Status: ", apiResp.Status)
		fmt.Println("")

	} else {
		apiResp := addQIErrorResponse{}
		jsonErr := json.Unmarshal(body, &apiResp)
		if jsonErr != nil {
			return jsonErr
		}

		fmt.Println("Error adding queue item")
		fmt.Println("Message: " + apiResp.Message)
		fmt.Println("Error Code: " + strconv.Itoa(apiResp.ErrorCode))
		fmt.Println("")

	}

	return nil

}
