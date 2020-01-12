package gonagios

import (
	"encoding/json"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Client used to store info required to communicate with Nagios
type Client struct {
	URL        string
	Token      string
	httpClient *http.Client
}

// NewClient creates a pointer to the client that will be used to send requests to Nagios
func NewClient(url, token string) *Client {
	httpClient := &http.Client{
		Timeout: time.Second * 5,
	}

	nagiosClient := &Client{
		URL:        url,
		Token:      token,
		httpClient: httpClient,
	}

	return nagiosClient
}

func (client *Client) sendRequest(httpRequest *http.Request) ([]byte, error) {
	addRequestHeaders(httpRequest)

	response, err := client.httpClient.Do(httpRequest)

	// TODO: Need to validate that when Nagios is unavailable, this err check will catch it
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := readAPIResponse(response.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}

// buildURL creates the URL to communicate with the Nagios API
// Depending on object type and HTTP method, the URL will vary
func (client *Client) buildURL(apiType, objectType, methodType string, objectInfo ...string) string {
	var nagiosURL strings.Builder

	nagiosURL.WriteString(client.URL)

	if strings.HasSuffix(client.URL, "/") {
		nagiosURL.WriteString("api/v1/")
	} else {
		nagiosURL.WriteString("/api/v1/")
	}

	// config or system for API type
	nagiosURL.WriteString(apiType + "/")

	// Nagios object type (host, service, etc)
	nagiosURL.WriteString(objectType + "/")

	// For PUT and DELETE, we have to do a few extra things
	if methodType == http.MethodPut || methodType == http.MethodDelete {
		if objectInfo != nil {
			// Append the object's name to the URL
			nagiosURL.WriteString(objectInfo[0])
			if objectType == "service" { // If it's a service, we need to tack on the service description
				nagiosURL.WriteString("/" + objectInfo[1])
			}
		}
	}

	nagiosURL.WriteString("?apikey=" + client.Token)
	nagiosURL.WriteString("&pretty=1")

	return nagiosURL.String()
}

func addRequestHeaders(request *http.Request) {
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Accept", "/")

	return
}

func (client *Client) get(data *url.Values, resourceInfo interface{}, nagiosURL string) error {
	request, err := http.NewRequest(http.MethodGet, nagiosURL, strings.NewReader(data.Encode()))

	if err != nil {
		return err
	}

	body, err := client.sendRequest(request)

	if err != nil {
		return err
	}

	return json.Unmarshal(body, resourceInfo)
}

func (client *Client) post(data *url.Values, nagiosURL string) ([]byte, error) {
	request, err := http.NewRequest(http.MethodPost, nagiosURL, strings.NewReader(data.Encode()))

	if err != nil {
		return nil, err
	}

	body, err := client.sendRequest(request)

	if err != nil {
		return nil, err
	}

	return body, nil
}

func (client *Client) put(nagiosURL string) ([]byte, error) {
	if strings.Contains(nagiosURL, " ") {
		nagiosURL = strings.Replace(nagiosURL, " ", "%20", -1)
	}
	request, err := http.NewRequest(http.MethodPut, nagiosURL, nil)

	if err != nil {
		return nil, err
	}

	body, err := client.sendRequest(request)

	if err != nil {
		return nil, err
	}

	return body, nil
}

func (client *Client) delete(data *url.Values, nagiosURL string) ([]byte, error) {
	request, err := http.NewRequest(http.MethodDelete, nagiosURL, strings.NewReader(data.Encode()))

	if err != nil {
		return nil, err
	}

	body, err := client.sendRequest(request)

	if err != nil {
		return nil, err
	}

	return body, nil
}

func (client *Client) applyConfig() error {
	nagiosURL := client.buildURL("system", "applyconfig", http.MethodPost)

	data := &url.Values{}

	_, err := client.post(data, nagiosURL)

	if err != nil {
		return err
	}

	return nil
}

// Function maps the elements of a string array to a single string with each value separated by commas
// Nagios expects a list of values supplied in this format via URL encoding
func mapArrayToString(sourceArray []interface{}) string {
	var destString strings.Builder

	for i, sourceObject := range sourceArray {
		// If this is the first time looping through, set the destination object euqal to the first element in array
		if i == 0 {
			destString.WriteString(sourceObject.(string))
		} else { // More than one element in array. Append a comma first before we add the next item
			destString.WriteString(",")
			destString.WriteString(sourceObject.(string))
		}
	}

	return destString.String()
}

// Function takes any boolean value, converts to integer and returns in string format
func convertBoolToIntToString(sourceVal bool) string {
	if sourceVal {
		return "1"
	}
	return "0"
}

// setURLParams loops through a struct object and returns a set of URL parameters
func setURLParams(nagiosObject interface{}) *url.Values {
	values := reflect.ValueOf(nagiosObject)
	var urlParams = &url.Values{}
	var tag string

	// If we are passing in a pointer to a struct, we need to get the actual result of what the struct is pointing to
	if values.Kind() == reflect.Ptr {
		values = values.Elem()
	}

	for i := 0; i < values.NumField(); i++ {
		var outputString strings.Builder
		curType := values.Field(i).Type().String()
		tags := strings.Split(values.Type().Field(i).Tag.Get("json"), ",")

		for k := range tags {
			if tags[k] != "omitempty" {
				tag = tags[k]
				break
			}
		}

		if curType == "string" {
			if values.Field(i).Interface().(string) != "" {
				urlParams.Add(tag, values.Field(i).Interface().(string))
			}
		} else if curType == "[]interface {}" {
			if values.Field(i).Interface() != nil {
				for j, val := range values.Field(i).Interface().([]interface{}) {
					if j > 0 {
						outputString.WriteString(",")
					}

					outputString.WriteString(val.(string))
				}
				urlParams.Add(tag, outputString.String())
			}
		} else if curType == "int" {
			if strconv.Itoa(values.Field(i).Interface().(int)) != "" {
				// We need the value to be a string but first need to cast it as an integer if that is what the type is in the struct
				urlParams.Add(tag, strconv.Itoa(values.Field(i).Interface().(int)))
			}
		} else if curType == "map[string]interface {}" {
			if values.Field(i).Interface() != nil {
				// We need to loop through the map and grab the key and value for each line
				// The value is an interface, so we need to then call the Interface() method
				// and cast it as a string to get the value in string format
				mapObject := values.Field(i).MapRange()
				for mapObject.Next() {
					outputString.Reset()
					index := mapObject.Key().String()
					val := mapObject.Value()
					valString := val.Interface().(string)
					urlParams.Add(index, valString)
				}
			}
		}
	}
	return urlParams
}
