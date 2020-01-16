package gonagios

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createHostObject() *Host {
	hostName := "host1"
	hostAlias := "host1"
	hostAddress := "127.0.0.2"
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

	return host
}

func TestHost_newHost(t *testing.T) {
	if errList := envVarCheck(); errList != nil {
		t.Fatal(errList)
	}

	client := NewClient(os.Getenv("NAGIOS_URL"), os.Getenv("API_TOKEN"))

	host := createHostObject()

	body, err := client.NewHost(host)

	assert.NoError(t, err)

	responseCode := &ResponseCode{}

	err = json.Unmarshal(body, &responseCode)

	assert.NoError(t, err)
	assert.NotEmpty(t, responseCode.ResponseSuccess)
}

func TestHost_getHost(t *testing.T) {
	if errList := envVarCheck(); errList != nil {
		t.Fatal(errList)
	}

	hostName := "host1"

	client := NewClient(os.Getenv("NAGIOS_URL"), os.Getenv("API_TOKEN"))

	host, err := client.GetHost(hostName)

	assert.NoError(t, err)
	assert.NotEmpty(t, host)
	assert.Equal(t, "host1", host.HostName)
}

func TestHost_updateHost(t *testing.T) {
	if errList := envVarCheck(); errList != nil {
		t.Fatal(errList)
	}

	client := NewClient(os.Getenv("NAGIOS_URL"), os.Getenv("API_TOKEN"))
	// time.Sleep(10 * time.Second)
	host := createHostObject()

	oldName := "host1"
	host.HostName = "updatedHost1"

	// Update the host in Nagios
	err := client.UpdateHost(host, oldName)

	assert.NoError(t, err)

	// Get the new host from Nagios to ensure the update was successful
	updatedHost, err := client.GetHost(host.HostName)

	assert.Equal(t, "updatedHost1", updatedHost.HostName)
}

func TestHost_deleteHost(t *testing.T) {
	if errList := envVarCheck(); errList != nil {
		t.Fatal(errList)
	}

	hostName := "updatedHost1"

	client := NewClient(os.Getenv("NAGIOS_URL"), os.Getenv("API_TOKEN"))

	// Delete host from Nagios
	body, err := client.DeleteHost(hostName)

	assert.NoError(t, err)

	// Lets make sure that the response code that is returned is 'success'
	responseCode := &ResponseCode{}
	err = json.Unmarshal(body, &responseCode)

	assert.NoError(t, err)
	assert.NotEmpty(t, responseCode.ResponseSuccess)
}
