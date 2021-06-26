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

// IP is a structure to store details about ip addresses
type IP struct {
	gorm.Model
	ID         string       `json:"id" gorm:"PrimaryKey"`
	Subdomains []*Subdomain `json:"subdomains" gorm:"many2many:subdomain_ips;"`
	ProgramID  string       `json:"program"`
}

// GetIPs will get all ips from database
func (c *Client) GetIPs() ([]IP, error) {
	rel := &url.URL{Path: "/api/ips"}
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
	var ips []IP
	err = json.NewDecoder(resp.Body).Decode(&ips)
	return ips, err
}

// GetIP will get a ip
func (c *Client) GetIP(id string) (IP, error) {
	var emptyip IP
	rel := &url.URL{Path: "/api/ips/" + id}
	u := c.BaseURL.ResolveReference(rel)
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return emptyip, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return emptyip, err
	}
	defer resp.Body.Close()
	var ip IP
	err = json.NewDecoder(resp.Body).Decode(&ip)
	return ip, err
}

// CreateIPs will create new IPs
func (c *Client) CreateIPs(ips []IP) ([]IP, error) {
	var emptyip []IP

	jsonips, err := json.Marshal(ips)
	if err != nil {
		log.Println("Could not convert ip to JSON, is it in the correct format?")
	}
	rel := &url.URL{Path: "/api/ips"}
	u := c.BaseURL.ResolveReference(rel)
	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(jsonips))
	if err != nil {
		return emptyip, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return emptyip, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&ips)
	return ips, err
}

// UpdateIP will update the specified ip
func (c *Client) UpdateIP(ip IP, id string) (IP, error) {
	var emptyip IP
	jsonip, err := json.Marshal(ip)
	if err != nil {
		log.Println("Could not convert ip to JSON, is it in the correct format?")
	}
	rel := &url.URL{Path: "/api/ips/" + id}
	u := c.BaseURL.ResolveReference(rel)
	req, err := http.NewRequest("PUT", u.String(), bytes.NewBuffer(jsonip))
	if err != nil {
		return emptyip, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return emptyip, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&ip)
	return ip, err
}

// DeleteIP will delete a ip
func (c *Client) DeleteIP(id string) (bool, error) {
	rel := &url.URL{Path: "/api/ips/" + id}
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
	return true, err
}

// printIPs prints multiple ips to terminal in desired output format
func printIPs(outputFormat string, ips []IP) {

	// output in desired format
	if outputFormat == "json" {

		ipsJSON, err := json.Marshal(ips)
		if err != nil {
			fmt.Println("Error occured while converting the response to JSON: ", err)
		}
		fmt.Println(string(ipsJSON))

	} else {

		for _, ip := range ips {
			fmt.Println(ip.ID)
		}

	}
}

// printIP prints the ip to terminal in desired output format
func printIP(ipID string, outputFormat string, c Client) {
	// get the ip
	ip, err := c.GetIP(ipID)
	if err != nil {
		fmt.Println("Error occured while fetching ip.", err)
	}

	// output in desired format
	if outputFormat == "json" {

		ipJSON, err := json.Marshal(ip)
		if err != nil {
			fmt.Println("Error occured while converting the response to JSON: ", err)
		}
		fmt.Println(string(ipJSON))

	} else {
		fmt.Println(ip.ID)
	}
}

// IPsCLI handles the ips subcommand CLI
func IPsCLI(c Client) {
	if len(os.Args) < 3 {
		fmt.Println("Invalid arguments. Hint: ./hakstore-client ips {list|create|delete}")
		return
	}
	switch os.Args[2] {
	case "list":
		ipsFlagSet := flag.NewFlagSet("ips list", flag.ExitOnError)
		ipID := ipsFlagSet.String("id", "", "ID of ip")
		outputFormat := ipsFlagSet.String("output", "", "output format")
		programID := ipsFlagSet.String("program", "", "ID of program")
		ipsFlagSet.Parse(os.Args[3:])
		if isFlagPassed("id", ipsFlagSet) {
			// show single ip
			printIP(*ipID, *outputFormat, c)

		} else if isFlagPassed("program", ipsFlagSet) {
			PrintAssociatedIPs(*programID, *outputFormat, c)
		} else {
			// list ips
			ips, err := c.GetIPs()
			if err != nil {
				fmt.Println("Error retreiving ips: ", err)
			}
			printIPs(*outputFormat, ips)
		}
	case "create":
		ipsFlagSet := flag.NewFlagSet("ips create", flag.ExitOnError)
		ipID := ipsFlagSet.String("id", "", "ID of ip, e.g. api.example.com")
		programID := ipsFlagSet.String("program", "", "Program that the ip will be associated with, e.g. tesla")
		ipsFlagSet.Parse(os.Args[3:])
		// create/update ip
		if *ipID == "" || *programID == "" {
			fmt.Println("You need to specify a -id and a -program to create the ip.")
			return
		}
		ips := []IP{{ID: *ipID, ProgramID: *programID}}
		_, err := c.CreateIPs(ips)
		if err != nil {
			fmt.Println("An error occured while creating the ip: ", err)
		}
	case "delete":
		ipsFlagSet := flag.NewFlagSet("ips delete", flag.ExitOnError)
		ipID := ipsFlagSet.String("id", "", "ID of ip")
		ipsFlagSet.Parse(os.Args[3:])
		// delete ip
		if *ipID == "" {
			fmt.Println("You need to specify a ip id to delete with -id.")
			return
		}
		c.DeleteIP(*ipID)

	// no valid subcommand found - default to showing a message and exiting
	default:
		fmt.Println("Invalid subsubcommand, ./hakstore-client ips {list|create|delete}")
		os.Exit(1)
	}
}
