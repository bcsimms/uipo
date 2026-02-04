package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type httpErrorResp struct {
	Message string `json:"message"`
}

// HTTPHelper is used by all commands to send their http requests and
// handle invalid responses
func HTTPHelper(client http.Client, req *http.Request) ([]byte, error) {

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error communicating with the API endpoint")
		fmt.Println(err.Error())
		return nil, err
	}

	defer resp.Body.Close()
	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, readErr
	}

	if resp.StatusCode != 200 {
		fmt.Println("\nAPI Request Failed.  Here's the response:")
		fmt.Println("  ", resp.Status)
		errResp := httpErrorResp{}
		jsonErr := json.Unmarshal(body, &errResp)
		if jsonErr != nil {
			return nil, jsonErr
		}
		fmt.Println("  ", errResp.Message)
		return nil, errors.New("Request failed")
	}

	return body, nil

}
