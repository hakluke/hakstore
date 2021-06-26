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

// RootDomain is a structure to store details about bug bounty rootdomains
type RootDomain struct {
	gorm.Model
	ID         string      `json:"id" gorm:"PrimaryKey"`
	ProgramID  string      `json:"program"`
	Subdomains []Subdomain `json:"subdomains"`
}

// GetRootDomains will get all rootdomains from database
func (c *Client) GetRootDomains() ([]RootDomain, error) {
	rel := &url.URL{Path: "/api/rootdomains"}
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
	var rootdomains []RootDomain
	err = json.NewDecoder(resp.Body).Decode(&rootdomains)
	return rootdomains, err
}

// GetRootDomain will get a rootdomain
func (c *Client) GetRootDomain(id string) (RootDomain, error) {
	var emptyrootdomain RootDomain
	rel := &url.URL{Path: "/api/rootdomains/" + id}
	u := c.BaseURL.ResolveReference(rel)
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return emptyrootdomain, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return emptyrootdomain, err
	}
	defer resp.Body.Close()
	var rootdomain RootDomain
	err = json.NewDecoder(resp.Body).Decode(&rootdomain)
	return rootdomain, err
}

// CreateRootDomain will create a new rootdomain
func (c *Client) CreateRootDomain(rootdomain RootDomain) (RootDomain, error) {
	var emptyrootdomain RootDomain
	jsonrootdomain, err := json.Marshal(rootdomain)
	if err != nil {
		log.Println("Could not convert rootdomain to JSON, is it in the correct format?")
	}
	rel := &url.URL{Path: "/api/rootdomains"}
	u := c.BaseURL.ResolveReference(rel)
	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(jsonrootdomain))
	if err != nil {
		return emptyrootdomain, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return emptyrootdomain, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&rootdomain)
	return rootdomain, err
}

// UpdateRootDomain will update the specified rootdomain
func (c *Client) UpdateRootDomain(rootdomain RootDomain, id string) (RootDomain, error) {
	var emptyrootdomain RootDomain
	jsonrootdomain, err := json.Marshal(rootdomain)
	if err != nil {
		log.Println("Could not convert rootdomain to JSON, is it in the correct format?")
	}
	rel := &url.URL{Path: "/api/rootdomains/" + id}
	u := c.BaseURL.ResolveReference(rel)
	req, err := http.NewRequest("PUT", u.String(), bytes.NewBuffer(jsonrootdomain))
	if err != nil {
		return emptyrootdomain, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return emptyrootdomain, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&rootdomain)
	return rootdomain, err
}

// DeleteRootDomain will get a rootdomain
func (c *Client) DeleteRootDomain(id string) (bool, error) {
	rel := &url.URL{Path: "/api/rootdomains/" + id}
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
	var rootdomain RootDomain
	err = json.NewDecoder(resp.Body).Decode(&rootdomain)
	return true, err
}

// GetAssociatedSubdomains will get subdomains belonging to the specified rootdomain
func (c *Client) GetAssociatedSubdomains(id string) ([]Subdomain, error) {
	rel := &url.URL{Path: "/api/rootdomains/" + id + "/subdomains"}
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
	var subdomains []Subdomain
	err = json.NewDecoder(resp.Body).Decode(&subdomains)
	return subdomains, err
}

// PrintRootDomains prints the rootdomain to terminal in desired output format
func PrintRootDomains(outputFormat string, c Client) {
	// get the rootdomains
	rootdomains, err := c.GetRootDomains()
	if err != nil {
		fmt.Println("Error occured while fetching rootdomains.", err)
	}

	// output in desired format
	if outputFormat == "json" {

		rootdomainsJSON, err := json.Marshal(rootdomains)
		if err != nil {
			fmt.Println("Error occured while converting the response to JSON: ", err)
		}
		fmt.Println(string(rootdomainsJSON))

	} else {

		for _, rootdomain := range rootdomains {
			fmt.Println(rootdomain.ID)
		}

	}
}

// PrintRootDomain prints the rootdomain to terminal in desired output format
func PrintRootDomain(rootdomainID string, outputFormat string, c Client) {
	// get the rootdomain
	rootdomain, err := c.GetRootDomain(rootdomainID)
	if err != nil {
		fmt.Println("Error occured while fetching rootdomain.", err)
	}

	// output in desired format
	if outputFormat == "json" {

		rootdomainJSON, err := json.Marshal(rootdomain)
		if err != nil {
			fmt.Println("Error occured while converting the response to JSON: ", err)
		}
		fmt.Println(string(rootdomainJSON))

	} else {
		fmt.Println(rootdomain.ID)
	}
}

// PrintAssociatedSubdomains will print all rootdomains associated with the platform
func PrintAssociatedSubdomains(id string, outputFormat string, c Client) {
	subdomains, err := c.GetAssociatedSubdomains(id)
	if err != nil {
		fmt.Println("Error occured while fetching associated subdomains: ", err)
		return
	}
	// output in desired format
	if outputFormat == "json" {
		subdomainsJSON, err := json.Marshal(subdomains)
		if err != nil {
			fmt.Println("Error occured while converting the response to JSON: ", err)
		}
		fmt.Println(string(subdomainsJSON))
	} else {
		for _, subdomain := range subdomains {
			fmt.Println(subdomain.ID)
		}
	}
}

// RootdomainsCLI handles the rootdomain subcommand CLI
func RootdomainsCLI(c Client) {
	if len(os.Args) < 3 {
		fmt.Println("Invalid arguments. Hint: ./hakstore-client rootdomains {list|create|delete}")
		return
	}
	switch os.Args[2] {
	case "list":
		rootdomainsFlagSet := flag.NewFlagSet("rootdomains list", flag.ExitOnError)
		rootdomainID := rootdomainsFlagSet.String("id", "", "ID of rootdomain")
		outputFormat := rootdomainsFlagSet.String("output", "", "output format")
		programID := rootdomainsFlagSet.String("program", "", "ID of program")
		rootdomainsFlagSet.Parse(os.Args[3:])
		if isFlagPassed("id", rootdomainsFlagSet) {
			// show single rootdomain
			PrintRootDomain(*rootdomainID, *outputFormat, c)

		} else if isFlagPassed("program", rootdomainsFlagSet) {
			PrintAssociatedRootDomains(*programID, *outputFormat, c)
		} else {
			// list rootdomains
			PrintRootDomains(*outputFormat, c)
		}
	case "create":
		rootdomainsFlagSet := flag.NewFlagSet("rootdomains create", flag.ExitOnError)
		rootdomainID := rootdomainsFlagSet.String("id", "", "ID of rootdomain")
		program := rootdomainsFlagSet.String("program", "", "Program that rootdomain is associated with")
		rootdomainsFlagSet.Parse(os.Args[3:])
		// create/update rootdomain
		if *rootdomainID == "" || *program == "" {
			fmt.Println("You need to specify a -id and a -program to create the rootdomain.")
			return
		}
		rootdomain := RootDomain{ID: *rootdomainID, ProgramID: *program}
		_, err := c.CreateRootDomain(rootdomain)
		if err != nil {
			fmt.Println("An error occured while creating the rootdomain: ", err)
		}
	case "delete":
		rootdomainsFlagSet := flag.NewFlagSet("rootdomains delete", flag.ExitOnError)
		rootdomainID := rootdomainsFlagSet.String("id", "", "ID of rootdomain")
		rootdomainsFlagSet.Parse(os.Args[3:])
		// delete rootdomain
		if *rootdomainID == "" {
			fmt.Println("You need to specify a rootdomain id to delete with -id.")
			return
		}
		c.DeleteRootDomain(*rootdomainID)

	// no valid subcommand found - default to showing a message and exiting
	default:
		fmt.Println("Invalid subsubcommand, ./hakstore-client rootdomains {list|create|delete}")
		os.Exit(1)
	}
}
