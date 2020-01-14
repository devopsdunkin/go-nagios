package gonagios

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var requiredEnvVars = []string{
	"NAGIOS_URL",
	"API_TOKEN",
}

func TestNewClient(t *testing.T) {
	apiType := "system"
	objectType := "info"

	envVarCheck(t)

	client := NewClient(os.Getenv("NAGIOS_URL"), os.Getenv("API_TOKEN"))

	// Test connection to Nagios with client
	// Query system/status
	infoURL := client.buildURL(apiType, objectType, http.MethodGet)
	data := &url.Values{}

	body, err := client.get(data.Encode(), infoURL)

	assert.NoError(t, err)

	response := struct {
		Product string `json:"product"`
	}{}

	err = json.Unmarshal(body, &response)

	assert.NoError(t, err)
	assert.Equal(t, "nagiosxi", response.Product)
}

func envVarCheck(t *testing.T) {
	var errList []string
	for _, variable := range requiredEnvVars {
		if os.Getenv(variable) == "" {
			err := "\n" + variable + " environment variable must be set to continue tests"
			errList = append(errList, err)
		}
	}

	if errList != nil {
		t.Fatal(errList)
	}
}
