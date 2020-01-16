# GOnagios

[![CircleCI](https://circleci.com/gh/devopsdunkin/gonagios/tree/master.svg?style=svg)](https://circleci.com/gh/devopsdunkin/gonagios/tree/master)

Go client for managing endpoints in Nagios XI

## Installing

`go get -u github.com/devopsdunkin/gonagios`

## Example

```go
package main

import (
    "log"

    "github.com/devopsdunkin/gonagios"
)

func main() {
    url := "https://nagios.domain.local/nagiosxi"
    token := "token123"
    client := gonagios.NewClient(url, token)

    hostContacts := make([]interface{}, 1)
    hostTemplates := make([]interface{}, 1)
    hostContacts[0] = "nagiosadmin"
    hostTemplates[0] = "generic-host"

    host := &Host{
        HostName:             "host1",
        Alias:                "host1",
        Address:              "192.168.1.1",
        MaxCheckAttempts:     "5",
        CheckPeriod:          "24x7",
        NotificationInterval: "10",
        NotificationPeriod:   "24x7",
        Contacts:             hostContacts,
        Templates:            hostTemplates,
    }

    _, err := client.NewHost()

    if err != nil {
        log.Fatal(err)
    }

    _, err := client.DeleteHost(host.HostName)

    if err != nil {
        log.Fatal(err)
    }
}
```
