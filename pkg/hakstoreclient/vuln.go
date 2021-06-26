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

// Vuln is the model that holds details about a vulnerability. Currently IPs/Subdomains are a many to many relationship,
// but I think they will ultimately hold only one value per vuln.
type Vuln struct {
	gorm.Model
	//ID         string       `json:"id" gorm:"PrimaryKey;autoIncrement"`
	ID          int          `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Subdomains  []*Subdomain `json:"subdomains" gorm:"many2many:subdomain_vulns;"`
	IPs         []*IP        `json:"vulns" gorm:"many2many:subdomain_vulns;"`
	Description string       `json:"description"`
	ProgramID   string       `json:"program"`
	Severity    int          `json:"severity"`
}

// GetVulns will get all vulns from database
func (c *Client) GetVulns() ([]Vuln, error) {
	rel := &url.URL{Path: "/api/vulns"}
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
	var vulns []Vuln
	err = json.NewDecoder(resp.Body).Decode(&vulns)
	return vulns, err
}

// GetVuln will get a vuln
func (c *Client) GetVuln(id string) (Vuln, error) {
	var emptyvuln Vuln
	rel := &url.URL{Path: "/api/vulns/" + id}
	u := c.BaseURL.ResolveReference(rel)
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return emptyvuln, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return emptyvuln, err
	}
	defer resp.Body.Close()
	var vuln Vuln
	err = json.NewDecoder(resp.Body).Decode(&vuln)
	return vuln, err
}

// CreateVulns will create new Vulns
func (c *Client) CreateVulns(vulns []Vuln) ([]Vuln, error) {
	var emptyvuln []Vuln

	jsonvulns, err := json.Marshal(vulns)
	if err != nil {
		log.Println("Could not convert vuln to JSON, is it in the correct format?")
	}
	rel := &url.URL{Path: "/api/vulns"}
	u := c.BaseURL.ResolveReference(rel)
	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(jsonvulns))
	if err != nil {
		return emptyvuln, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return emptyvuln, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&vulns)
	return vulns, err
}

// UpdateVuln will update the specified vuln
func (c *Client) UpdateVuln(vuln Vuln, id string) (Vuln, error) {
	var emptyvuln Vuln
	jsonvuln, err := json.Marshal(vuln)
	if err != nil {
		log.Println("Could not convert vuln to JSON, is it in the correct format?")
	}
	rel := &url.URL{Path: "/api/vulns/" + id}
	u := c.BaseURL.ResolveReference(rel)
	req, err := http.NewRequest("PUT", u.String(), bytes.NewBuffer(jsonvuln))
	if err != nil {
		return emptyvuln, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return emptyvuln, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&vuln)
	return vuln, err
}

// DeleteVuln will delete a vuln
func (c *Client) DeleteVuln(id string) (bool, error) {
	rel := &url.URL{Path: "/api/vulns/" + id}
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

// printVulns prints multvulnle vulns to terminal in desired output format
func printVulns(outputFormat string, vulns []Vuln) {

	// output in desired format
	if outputFormat == "json" {

		vulnsJSON, err := json.Marshal(vulns)
		if err != nil {
			fmt.Println("Error occured while converting the response to JSON: ", err)
		}
		fmt.Println(string(vulnsJSON))

	} else {

		for _, vuln := range vulns {
			fmt.Println(vuln.ID)
		}

	}
}

// printVuln prints the vuln to terminal in desired output format
func printVuln(vulnID string, outputFormat string, c Client) {
	// get the vuln
	vuln, err := c.GetVuln(vulnID)
	if err != nil {
		fmt.Println("Error occured while fetching vuln.", err)
	}

	// output in desired format
	if outputFormat == "json" {

		vulnJSON, err := json.Marshal(vuln)
		if err != nil {
			fmt.Println("Error occured while converting the response to JSON: ", err)
		}
		fmt.Println(string(vulnJSON))

	} else {
		fmt.Println(vuln.ID)
	}
}

// VulnsCLI handles the vulns subcommand CLI
func VulnsCLI(c Client) {
	if len(os.Args) < 3 {
		fmt.Println("Invalid arguments. Hint: ./hakstore-client vulns {list|create|delete}")
		return
	}
	switch os.Args[2] {
	case "list":
		vulnsFlagSet := flag.NewFlagSet("vulns list", flag.ExitOnError)
		vulnID := vulnsFlagSet.String("id", "", "ID of vuln")
		outputFormat := vulnsFlagSet.String("output", "", "output format")
		programID := vulnsFlagSet.String("program", "", "ID of program")

		vulnsFlagSet.Parse(os.Args[3:])
		if isFlagPassed("id", vulnsFlagSet) {
			// show single vuln
			printVuln(*vulnID, *outputFormat, c)

		} else if isFlagPassed("program", vulnsFlagSet) {
			PrintAssociatedVulns(*programID, *outputFormat, c)
		} else {
			// list vulns
			vulns, err := c.GetVulns()
			if err != nil {
				fmt.Println("Error retreiving vulns: ", err)
			}
			printVulns(*outputFormat, vulns)
		}
	case "create":
		vulnsFlagSet := flag.NewFlagSet("vulns create", flag.ExitOnError)
		description := vulnsFlagSet.String("description", "", "Description of vulnerability.")
		programID := vulnsFlagSet.String("program", "", "Program that the vuln will be associated with, e.g. tesla")
		severity := vulnsFlagSet.Int("severity", 1, "Severity of vulnerability from 1-5, 1 is critical, 5 is informational.")
		subdomain := vulnsFlagSet.String("subdomain", "", "Subdomain that vuln is associated with")
		ip := vulnsFlagSet.String("ip", "", "IP that vuln is associated with")

		vulnsFlagSet.Parse(os.Args[3:])
		// create/update vuln
		if *description == "" || *programID == "" || *severity == 0 {
			fmt.Println("You need to specify a -description, -program, -severity and -subdomain or -ip to create the vuln.")
			return
		}

		if *subdomain == "" && *ip == "" {
			fmt.Println("You need to specify either a -subdomain or an -ip to associate the vuln with")
		}

		vulns := []Vuln{{
			Description: *description,
			ProgramID:   *programID,
			Severity:    *severity,
			Subdomains: []*Subdomain{
				{ID: *subdomain},
			},
			IPs: []*IP{
				{ID: *ip},
			},
		}}
		_, err := c.CreateVulns(vulns)
		if err != nil {
			fmt.Println("An error occured while creating the vuln: ", err)
		}
	case "delete":
		vulnsFlagSet := flag.NewFlagSet("vulns delete", flag.ExitOnError)
		vulnID := vulnsFlagSet.String("id", "", "ID of vuln")
		vulnsFlagSet.Parse(os.Args[3:])
		// delete vuln
		if *vulnID == "" {
			fmt.Println("You need to specify a vuln id to delete with -id.")
			return
		}
		c.DeleteVuln(*vulnID)

	// no valid subcommand found - default to showing a message and exiting
	default:
		fmt.Println("Invalid subsubcommand, ./hakstore-client vulns {list|create|delete}")
		os.Exit(1)
	}
}
