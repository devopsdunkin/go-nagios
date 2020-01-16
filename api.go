package gonagios

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
)

// ResponseCode contains the list of available responses from Nagios
// It is used grab messages from Nagios to determine if the action
// was successful or an error occurred
type ResponseCode struct {
	ResponseSuccess string `json:"success"`
	ReponseError    string `json:"error"`
}

// readAPIResponse reads the body of the HTTP response sent from the Nagios XI API
// Nagios does not return errors in a way that golang will catch them in the err variable
// We need to enhance the io.ReadAll function with determing if Nagios returns a response of
// 'success' or 'error'
func parseAPIResponse(reader io.Reader) ([]byte, error) {
	var errorCode error
	body, _ := ioutil.ReadAll(reader)

	responseCode := &ResponseCode{}

	err := json.Unmarshal(body, &responseCode)

	if err != nil {
		return nil, err
	}

	if responseCode.ReponseError != "" {
		errorCode = errors.New(responseCode.ReponseError)
	}

	return body, errorCode
}
