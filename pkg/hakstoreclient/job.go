package hakstoreclient

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

// Job is a structure to store details about jobs that get sent to redis
type Job struct {
	Queue  string `json:"queue"`
	Target string `json:"target"`
}

// CreateJobs will create jobs for workers
func (c *Client) CreateJobs(jobs []Job) Message {
	jsonjobs, err := json.Marshal(jobs)
	if err != nil {
		log.Println("Could not convert jobs to JSON, is it in the correct format?")
	}
	rel := &url.URL{Path: "/api/jobs"}
	u := c.BaseURL.ResolveReference(rel)
	req, _ := http.NewRequest("POST", u.String(), bytes.NewBuffer(jsonjobs))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return Message{Success: false, Message: "Failed to connect to server."}
	}
	defer resp.Body.Close()
	var message Message
	err = json.NewDecoder(resp.Body).Decode(&message)
	return message
}

// JobsCLI handles jobs that are sent to redis for processing by workers
func JobsCLI(c Client) {
	if len(os.Args) < 3 {
		fmt.Println("Invalid arguments. Hint: ./hakstore-client jobs {list|create|delete}")
		return
	}

	switch os.Args[2] {
	case "create":
		jobsFlagSet := flag.NewFlagSet("jobs list", flag.ExitOnError)
		queue := jobsFlagSet.String("queue", "", "job queue")
		target := jobsFlagSet.String("target", "", "job target")
		jobsFlagSet.Parse(os.Args[3:])

		if !isFlagPassed("queue", jobsFlagSet) || !isFlagPassed("target", jobsFlagSet) {
			// show single subdomain
			fmt.Println("You need to provide a -queue and -target. Queue is the queue that the job will be pushed to (for example, nuclei or updateDNSData) and target is the actual subdomain that the task will be performed against.")
		}

		// if "all" is specified as the target, get all subs from the database and add them all, otherwise just add the specified one
		if *target == "all" {
			subs, err := c.GetSubdomains()
			if err != nil {
				log.Fatal("Failed to get subdomains.")
			}
			var jobs []Job
			for _, sub := range subs {
				jobs = append(jobs, Job{Queue: *queue, Target: sub.ID})
			}
			message := c.CreateJobs(jobs)
			fmt.Println(message.Message)
		} else {
			jobs := []Job{{Queue: *queue, Target: *target}}
			message := c.CreateJobs(jobs)
			fmt.Println(message.Message)
		}

	default:
		fmt.Println("Invalid subsubcommand, ./hakstore-client jobs create ...")
		os.Exit(1)
	}
}
