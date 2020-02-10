package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/bcsimms/uipo/config"
	"github.com/bcsimms/uipo/util"
)

// CmdAuthenticate is the struct that represents the command line options for the "authenticate" command
//   UserID and RefreshToken are mutually exclusive and used to determine the authentication mode
//   Two modes supported:
//		1)  On-Premise - using basic authentication to retrieve a bearer token
//		2)  UiPath Platform (SaaS) - using a refresh token to retrieve a bearer token
//   Either mode results in a Bearer Token used for subsequent calls to the UiPath Orchestrator
// Some data needed for authentication will be cahced in the .uipo/config.json file for streamlined authentication requests
//   This includes:
//     - Authentication endpoint
//     - Refresh token
//     - [add more here as we go]
type CmdAuthenticate struct {
	AuthorizationEndpoint string `short:"e" long:"endpoint" description:"API endpoint (e.g. https://api.example.com)"`
	Tenant                string `short:"t" long:"tenant" description:"Tenant"`
	UserID                string `short:"u" long:"userid" description:"UserID - Used for on-premise installations"`
	Password              string `short:"p" long:"password" description:"Password - Used for on-premise installations"`
	RefreshToken          string `short:"r" long:"refreh-token" description:"Refresh Token - Used for UiPath Platform Installations"`
	ClientID              string `short:"c" long:"client-id" default:"5v7PmPJL6FOGu6RB8I1Y4adLBhIwovQN" descritpion:"Client ID - Used for UiPath Platform Installations.  Should not need to be overridden"`

	Config Config
}

// Used when making authentication calls into a UiPath hosted Orchestrator
type requestBodyHosted struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	RefreshToken string `json:"refresh_token"`
}

// Used when making authentication calls into a Orchestrator hosted on-premise
type requestBodyOnPrem struct {
	TenantName string `json:"tenancyName"`
	UserID     string `json:"usernameOrEmailAddress"`
	Password   string `json:"password"`
}

// Return structure for a hosted authentication call
type responseHosted struct {
	AccessToken string `json:"access_token"`
	IDToken     string `json:"id_token"`
	Scope       string `json:"scope"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type responseOnPrem struct {
	AccessToken         string `json:"result"`
	TargetURL           string `json:"targetUrl"`
	Success             bool   `json:"success"`
	Error               string `json:"error"`
	UnAuthorizedRequest bool   `json:"unAuthorizedRequest"`
	ABP                 bool   `json:"__abp"`
}

// The ahtentication type, determined by the arugments provided
type authType int

const (
	authInvalid authType = iota
	authOnPrem
	authHosted
)

// Setup is the override of the ExtendedCommander
func (cmd *CmdAuthenticate) Setup(config Config) error {

	util.LogDebug("Starting Authenticate Setup")
	cmd.Config = config

	//Setup cached values from config file
	if cmd.Config.GetAuthorizationEndpoint() != "" {
		cmd.AuthorizationEndpoint = cmd.Config.GetAuthorizationEndpoint()
	}
	if cmd.Config.GetRefreshToken() != "" {
		cmd.RefreshToken = cmd.Config.GetRefreshToken()
	}

	return nil
}

// Execute is the override for the standard command execute method
//  This method will perform the sequence of steps to authenticate a user
func (cmd *CmdAuthenticate) Execute(args []string) error {

	util.LogDebug("Starting Authenticate Execute")
	authMode, err := cmd.validateFlags()
	if err != nil {
		return err
	}

	if authMode == authHosted {

		util.LogDebug("Authentication Mode Hosted")

		cmd.Config.SetEndpointType(config.EndpointTypeHosted)

		reqBody := requestBodyHosted{}
		reqBody.GrantType = "refresh_token"
		reqBody.ClientID = cmd.ClientID
		reqBody.RefreshToken = cmd.RefreshToken

		requestBody, err := json.Marshal(&reqBody)
		if err != nil {
			return err
		}

		// Send the authentication request
		util.LogDebug("Sending authentication request")
		resp, err := http.Post(cmd.AuthorizationEndpoint, "application/json", strings.NewReader(string(requestBody)))

		if err != nil {
			return err
		}

		util.LogDebug("Reading response body")
		body, readErr := ioutil.ReadAll(resp.Body)
		util.LogDebug("Response was: \n" + string(body))
		if readErr != nil {
			return readErr
		}
		defer resp.Body.Close()

		oauthResp := responseHosted{}

		util.LogDebug("Reading JSon response data")
		jsonErr := json.Unmarshal(body, &oauthResp)
		if jsonErr != nil {
			return jsonErr
		}

		// Save our Access Token (Bearer token) to cache in the config file
		cmd.Config.SetAccessToken(oauthResp.AccessToken)

		if cmd.Config.GetAuthorizationEndpoint() == "" {
			cmd.Config.SetAuthorizationEndpoint(cmd.AuthorizationEndpoint)
		}

		fmt.Println("Authentication successful.  Bearer token cached for future requests.")

		expiryDateTime := time.Now().Add(time.Second * time.Duration(oauthResp.ExpiresIn))

		fmt.Println("Expires on: ", expiryDateTime)

	} else if authMode == authOnPrem {

		util.LogDebug("Authentication Mode On-Premise")

		cmd.Config.SetEndpointType(config.EndpointTypeOnPremise)

		reqBody := requestBodyOnPrem{}
		reqBody.TenantName = cmd.Tenant
		reqBody.UserID = cmd.UserID
		reqBody.Password = cmd.Password

		requestBody, err := json.Marshal(&reqBody)
		if err != nil {
			return err
		}

		util.LogDebug("Sending authentication request")
		// Send the authentication request
		resp, err := http.Post(cmd.AuthorizationEndpoint, "application/json", strings.NewReader(string(requestBody)))

		if err != nil {
			return err
		}

		util.LogDebug("Reading response body")
		body, readErr := ioutil.ReadAll(resp.Body)
		util.LogDebug("Response was: \n" + string(body))
		if readErr != nil {
			return readErr
		}
		defer resp.Body.Close()

		//fmt.Println(string(body))
		oauthResp := responseOnPrem{}

		util.LogDebug("Reading JSon response data")
		jsonErr := json.Unmarshal(body, &oauthResp)
		if jsonErr != nil {
			return jsonErr
		}

		// Save our Access Token (Bearer token) to cache in the config file
		cmd.Config.SetAccessToken(oauthResp.AccessToken)

		if cmd.Config.GetRefreshToken() == "" {
			cmd.Config.SetRefreshToken(cmd.RefreshToken)
		}

		if cmd.Config.GetAuthorizationEndpoint() == "" {
			cmd.Config.SetAuthorizationEndpoint(cmd.AuthorizationEndpoint)
		}

		fmt.Println("Authentication successful.  Bearer token cached for future requests.")

	}

	return nil
}

func (cmd *CmdAuthenticate) validateFlags() (authType, error) {

	// If a UserID is provided on the command line, use on-premise as the default auth method
	// This will allow an end user to also setup, and cache, hosted Orchestrator settings and
	//  will allow the end user to use a single configuration for either hosted or on-premise
	//  installations
	// If we are using an on-premise Ochestrator, make sure we have the flags we need
	if cmd.UserID != "" {

		if cmd.Password == "" || cmd.Tenant == "" {
			return authInvalid, errors.New("When using on-premise authentication, password and tenant are required arguments")
		}
		return authOnPrem, nil
	}

	// If we are using a hosted Orchestrator, make sure we have the flags we need
	if cmd.RefreshToken != "" {

		//if cmd.AccountLogicalName == "" || cmd.ServiceLogicalName == "" {
		//	return authInvalid, errors.New("When using hosted authentication, Account Logical Name and Service Logical Name are required arguments")
		//} else {
		return authHosted, nil
		//}
	}

	return authInvalid, errors.New("No valid authentication type was able to be determined")
}
