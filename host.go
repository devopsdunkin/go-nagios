package gonagios

import (
	"net/http"
)

var apiType = "config"
var objectType = "host"

// Host contains all available attributes for a Nagios host object
type Host struct {
	HostName                   string                 `json:"host_name"`
	Address                    string                 `json:"address"`
	DisplayName                string                 `json:"display_name,omitempty"`
	MaxCheckAttempts           string                 `json:"max_check_attempts"`
	CheckPeriod                string                 `json:"check_period"`
	NotificationInterval       string                 `json:"notification_interval"`
	NotificationPeriod         string                 `json:"notification_period"`
	Contacts                   []interface{}          `json:"contacts"`
	Alias                      string                 `json:"alias,omitempty"`
	Templates                  []interface{}          `json:"use,omitempty"`
	CheckCommand               string                 `json:"check_command,omitempty"`
	ContactGroups              []interface{}          `json:"contact_groups,omitempty"`
	Notes                      string                 `json:"notes,omitempty"`
	NotesURL                   string                 `json:"notes_url,omitempty"`
	ActionURL                  string                 `json:"action_url,omitempty"`
	InitialState               string                 `json:"initial_state,omitempty"`
	RetryInterval              string                 `json:"retry_interval,omitempty"`
	PassiveChecksEnabled       string                 `json:"passive_checks_enabled,omitempty"`
	ActiveChecksEnabled        string                 `json:"active_checks_enabled,omitempty"`
	ObsessOverHost             string                 `json:"obsess_over_host,omitempty"`
	EventHandler               string                 `json:"event_handler,omitempty"`
	EventHandlerEnabled        string                 `json:"event_handler_enabled,omitempty"`
	FlapDetectionEnabled       string                 `json:"flap_detection_enabled,omitempty"`
	FlapDetectionOptions       []interface{}          `json:"flap_detection_options,omitempty"`
	LowFlapThreshold           string                 `json:"low_flap_threshold,omitempty"`
	HighFlapThreshold          string                 `json:"high_flap_threshold,omitempty"`
	ProcessPerfData            string                 `json:"process_perf_data,omitempty"`
	RetainStatusInformation    string                 `json:"retain_status_information,omitempty"`
	RetainNonstatusInformation string                 `json:"retain_nonstatus_information,omitempty"`
	CheckFreshness             string                 `json:"check_freshness,omitempty"`
	FreshnessThreshold         string                 `json:"freshness_threshold,omitempty"`
	FirstNotificationDelay     string                 `json:"first_notification_delay,omitempty"`
	NotificationOptions        string                 `json:"notification_options,omitempty"`
	NotificationsEnabled       string                 `json:"notifications_enabled,omitempty"`
	StalkingOptions            string                 `json:"stalking_options,omitempty"`
	IconImage                  string                 `json:"icon_image,omitempty"`
	IconImageAlt               string                 `json:"icon_image_alt,omitempty"`
	VRMLImage                  string                 `json:"vrml_image,omitempty"`
	StatusMapImage             string                 `json:"statusmap_image,omitempty"`
	TwoDCoords                 string                 `json:"2d_coords,omitempty"`
	ThreeDCoords               string                 `json:"3d_coords,omitempty"`
	Register                   string                 `json:"register,omitempty"`
	FreeVariables              map[string]interface{} `json:"free_variables,omitempty"`
}

// NewHost creates a host object in Nagios
func (client *Client) NewHost(host *Host) ([]byte, error) {
	url := client.buildURL(apiType, objectType, http.MethodPost)

	data := setURLParams(host)

	body, err := client.post(data, url)

	if err != nil {
		return nil, err
	}

	err = client.applyConfig()

	if err != nil {
		return nil, err
	}

	return body, nil
}

// GetHost retrieves an existing host from Nagios
func (client *Client) GetHost(name string) (*Host, error) {
	var hostArray = []Host{} // TODO: Is there a better way to parse JSON without having to use an array of Host objects?
	var host Host

	url := client.buildURL(apiType, objectType, http.MethodGet)

	data := &url.Values{}
	data.Set("host_name", name)

	body, err := client.get(data.Encode(), )

	return nil, nil
}

// UpdateHost updates attributes of an existing host in Nagios
func UpdateHost() {
	return
}

// DeleteHost deletes a host from Nagios
func DeleteHost() {
	return
}
