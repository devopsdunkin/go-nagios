package gonagios

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHost_newHost(t *testing.T) {
	if errList := envVarCheck(); errList != nil {
		t.Fatal(errList)
	}

	hostName := "host1"
	hostAlias := "host1"
	hostAddress := "127.0.0.1"
	hostMaxCheckAttempts := "5"
	hostCheckPeriod := "24x7"
	hostNotificationInterval := "10"
	hostNotificationPeriod := "24x7"
	contact := "nagiosadmin"
	hostContacts := make([]interface{}, 1)
	hostContacts[0] = contact
	template := "generic-host"
	hostTemplates := make([]interface{}, 1)
	hostTemplates[0] = template
	client := NewClient(os.Getenv("NAGIOS_URL"), os.Getenv("API_TOKEN"))

	host := &Host{
		HostName:             hostName,
		Alias:                hostAlias,
		Address:              hostAddress,
		MaxCheckAttempts:     hostMaxCheckAttempts,
		CheckPeriod:          hostCheckPeriod,
		NotificationInterval: hostNotificationInterval,
		NotificationPeriod:   hostNotificationPeriod,
		Contacts:             hostContacts,
		Templates:            hostTemplates,
	}

	body, err := client.NewHost(host)

	assert.NoError(t, err)

	responseCode := &ResponseCode{}

	err = json.Unmarshal(body, &responseCode)

	assert.NoError(t, err)
	assert.NotEmpty(t, responseCode.ResponseSuccess)
}
