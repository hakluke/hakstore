// currently this is just the example from https://github.com/go-co-op/gocron - going to update to do proper scheduled tasks

package main

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
)

func task() {
	fmt.Println("I am running task.")
}

func taskWithParams(a int, b string) {
	fmt.Println(a, b)
}

func example() {
	location, err := time.LoadLocation(config.Database.TimeZone)
	if err != nil {
		print("Error occured, most likely due to an invalid location in the config file:", err)
	}

	s := gocron.NewScheduler(location)
	s.Every(1).Day().At("19:46:00").Do(task)
	s.StartAsync()
}
