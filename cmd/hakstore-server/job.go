package main

import (
	"encoding/json"
	"net/http"
)

// Job is a structure to store details about jobs that get sent to redis
type Job struct {
	Queue  string `json:"queue"`
	Target string `json:"target"`
}

// createJob creates a job and sends it to redis - the queue is the job identifier and the target is the host it will be performed on (i.e. subdomain or rootdomain string)
func createJobLocal(queue string, target string) {
	redisClient.LPush(queue, target)
}

// this is an API endpoint to create jobs
func createJobs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var jobs []Job
	_ = json.NewDecoder(r.Body).Decode(&jobs)
	for _, job := range jobs {
		createJobLocal(job.Queue, job.Target)
	}
	response := Message{Success: true, Message: "Jobs created successfully."}
	json.NewEncoder(w).Encode(response)
}
