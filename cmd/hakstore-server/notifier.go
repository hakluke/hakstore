package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"
)

type slackRequestBody struct {
	Text string `json:"text"`
}

// SendVulnNotification takes a vulnerability and sends a notification to the appropriate slack channel
func SendVulnNotification(v Vuln) {
	var webhook string
	switch v.Severity {
	case 1:
		webhook = config.Slack.CriticalWebhook
	case 2:
		webhook = config.Slack.HighWebhook
	case 3:
		webhook = config.Slack.MediumWebhook
	case 4:
		webhook = config.Slack.LowWebhook
	case 5:
		webhook = config.Slack.InformationalWebhook
	}

	// make a string of subdomains associated with this vuln
	var subsString string
	length := len(v.Subdomains)
	for n, s := range v.Subdomains {
		subsString = subsString + s.ID
		if n != length-1 {
			subsString = subsString + " "
		}
	}

	// make a string of ips associated with this vuln
	var ipsString string
	length = len(v.IPs)
	for n, s := range v.IPs {
		ipsString = ipsString + s.ID
		if n != length-1 {
			ipsString = ipsString + " "
		}
	}

	// Send a slack notification about the vuln
	SendNotification(webhook, "["+strings.ToUpper(severityString(v.Severity))+"]\nProgram: "+v.ProgramID+"\nHosts: [ "+subsString+" "+ipsString+"]\nDescription: "+v.Description)
}

// SendNotification sends a message to the slack webhook in the config.yml file
func SendNotification(webhook string, content string) {
	err := SendSlackNotification(webhook, content)
	if err != nil {
		log.Fatal(err)
	}
}

// SendSlackNotification will send a message via slack to a webhook
func SendSlackNotification(webhookURL string, msg string) error {

	slackBody, _ := json.Marshal(slackRequestBody{Text: msg})
	req, err := http.NewRequest(http.MethodPost, webhookURL, bytes.NewBuffer(slackBody))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	if buf.String() != "ok" {
		return errors.New("Error response returned from Slack")
	}
	return nil
}
