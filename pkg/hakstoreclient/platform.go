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

	"gorm.io/gorm"
)

// Platform Struct (Model)
type Platform struct {
	gorm.Model
	ID       string    `json:"id" gorm:"PrimaryKey"`
	URL      string    `json:"url"`
	Programs []Program `json:"programs"`
}

// GetPlatforms will get all platforms
func (c *Client) GetPlatforms() ([]Platform, error) {
	rel := &url.URL{Path: "/api/platforms"}
	u := c.BaseURL.ResolveReference(rel)
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var platforms []Platform
	err = json.NewDecoder(resp.Body).Decode(&platforms)
	return platforms, err
}

// GetPlatform will get a platform
func (c *Client) GetPlatform(id string) (Platform, error) {
	var emptyplatform Platform
	rel := &url.URL{Path: "/api/platforms/" + id}
	u := c.BaseURL.ResolveReference(rel)
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return emptyplatform, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return emptyplatform, err
	}
	defer resp.Body.Close()
	var platform Platform
	err = json.NewDecoder(resp.Body).Decode(&platform)
	return platform, err
}

// CreatePlatform will create a new platform
func (c *Client) CreatePlatform(platform Platform) (Platform, error) {
	var emptyplatform Platform
	jsonplatform, err := json.Marshal(platform)
	if err != nil {
		log.Println("Could not convert platform to JSON, is it in the correct format?")
	}
	rel := &url.URL{Path: "/api/platforms"}
	u := c.BaseURL.ResolveReference(rel)
	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(jsonplatform))
	if err != nil {
		return emptyplatform, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return emptyplatform, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&platform)
	return platform, err
}

// UpdatePlatform will update the specified platform
func (c *Client) UpdatePlatform(platform Platform, id string) (Platform, error) {
	var emptyplatform Platform
	jsonplatform, err := json.Marshal(platform)
	if err != nil {
		log.Println("Could not convert platform to JSON, is it in the correct format?")
	}
	rel := &url.URL{Path: "/api/platforms/" + id}
	u := c.BaseURL.ResolveReference(rel)
	req, err := http.NewRequest("PUT", u.String(), bytes.NewBuffer(jsonplatform))
	if err != nil {
		return emptyplatform, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return emptyplatform, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&platform)
	return platform, err
}

// DeletePlatform will get a platform
func (c *Client) DeletePlatform(id string) (bool, error) {
	rel := &url.URL{Path: "/api/platforms/" + id}
	u := c.BaseURL.ResolveReference(rel)
	req, err := http.NewRequest("DELETE", u.String(), nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	var platform Platform
	err = json.NewDecoder(resp.Body).Decode(&platform)
	return true, err
}

// GetAssociatedPrograms will get programs associated with a platform
func (c *Client) GetAssociatedPrograms(id string) ([]Program, error) {
	rel := &url.URL{Path: "/api/platforms/" + id + "/programs"}
	u := c.BaseURL.ResolveReference(rel)
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var programs []Program
	err = json.NewDecoder(resp.Body).Decode(&programs)
	return programs, err
}

// PrintAssociatedPrograms will print all programs associated with the platform
func PrintAssociatedPrograms(id string, outputFormat string, c Client) {
	programs, err := c.GetAssociatedPrograms(id)
	if err != nil {
		fmt.Println("Error occured while fetching associated programs: ", err)
		return
	}
	// output in desired format
	if outputFormat == "json" {
		programsJSON, err := json.Marshal(programs)
		if err != nil {
			fmt.Println("Error occured while converting the response to JSON: ", err)
		}
		fmt.Println(string(programsJSON))
	} else {
		for _, program := range programs {
			fmt.Println(program.ID)
		}
	}
}

// PrintPlatforms prints the platform to terminal in desired output format
func PrintPlatforms(outputFormat string, c Client) {
	// get the platforms
	platforms, err := c.GetPlatforms()
	if err != nil {
		fmt.Println("Error occured while fetching platforms.", err)
		return
	}

	// output in desired format
	if outputFormat == "json" {

		platformsJSON, err := json.Marshal(platforms)
		if err != nil {
			fmt.Println("Error occured while converting the response to JSON: ", err)
		}
		fmt.Println(string(platformsJSON))

	} else {

		for _, platform := range platforms {
			fmt.Println(platform.ID, platform.URL)
		}

	}
}

// PrintPlatform prints the platform to terminal in desired output format
func PrintPlatform(platformID string, outputFormat string, c Client) {
	// get the platform
	platform, err := c.GetPlatform(platformID)
	if err != nil {
		fmt.Println("Error occured while fetching platform.", err)
	}

	// output in desired format
	if outputFormat == "json" {

		platformJSON, err := json.Marshal(platform)
		if err != nil {
			fmt.Println("Error occured while converting the response to JSON: ", err)
		}
		fmt.Println(string(platformJSON))

	} else {

		fmt.Println(platform.ID, platform.URL)

	}
}

// PlatformCLI handles the platform subcommand CLI
func PlatformCLI(c Client) {

	if len(os.Args) < 3 {
		fmt.Println("Invalid arguments. Hint: ./hakstore-client platforms {list|create|delete}")
		return
	}
	switch os.Args[2] {
	case "list":
		platformsFlagSet := flag.NewFlagSet("platforms list", flag.ExitOnError)
		platformID := platformsFlagSet.String("id", "", "ID of platform")
		platformOutputFormat := platformsFlagSet.String("output", "", "output format")
		platformsFlagSet.Parse(os.Args[3:])
		if isFlagPassed("id", platformsFlagSet) {
			// show single platform
			PrintPlatform(*platformID, *platformOutputFormat, c)

		} else {
			// list platforms
			PrintPlatforms(*platformOutputFormat, c)
		}
	case "create":
		platformsFlagSet := flag.NewFlagSet("platforms create", flag.ExitOnError)
		platformID := platformsFlagSet.String("id", "", "ID of platform")
		platformURL := platformsFlagSet.String("url", "", "URL of platform")
		platformsFlagSet.Parse(os.Args[3:])
		// create/update platform
		if *platformID == "" || *platformURL == "" {
			fmt.Println("You need to specify a -id and -url to create the platform.")
			return
		}
		platform := Platform{ID: *platformID, URL: *platformURL}
		_, err := c.CreatePlatform(platform)
		if err != nil {
			fmt.Println("An error occured while creating the platform: ", err)
		}
	case "delete":
		platformsFlagSet := flag.NewFlagSet("platforms delete", flag.ExitOnError)
		platformID := platformsFlagSet.String("id", "", "ID of platform")
		platformsFlagSet.Parse(os.Args[3:])
		// delete platform
		if *platformID == "" {
			fmt.Println("You need to specify a platform id to delete with -id.")
			return
		}
		c.DeletePlatform(*platformID)

	// no valid subcommand found - default to showing a message and exiting
	default:
		fmt.Println("Invalid subsubcommand, ./hakstore-client platforms {list|create|delete}")
		os.Exit(1)
	}

}
