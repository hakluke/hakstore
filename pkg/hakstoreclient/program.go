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

// Program is a structure to store details about bug bounty programs
type Program struct {
	gorm.Model
	ID          string       `json:"id" gorm:"PrimaryKey"`
	PlatformID  string       `json:"platform"`
	RootDomains []RootDomain `json:"rootdomains"`
}

// GetPrograms will get all programs from database
func (c *Client) GetPrograms() ([]Program, error) {
	rel := &url.URL{Path: "/api/programs"}
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

// GetProgram will get a program
func (c *Client) GetProgram(id string) (Program, error) {
	var emptyprogram Program
	rel := &url.URL{Path: "/api/programs/" + id}
	u := c.BaseURL.ResolveReference(rel)
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return emptyprogram, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return emptyprogram, err
	}
	defer resp.Body.Close()
	var program Program
	err = json.NewDecoder(resp.Body).Decode(&program)
	return program, err
}

// CreateProgram will create a new program
func (c *Client) CreateProgram(program Program) (Program, error) {
	var emptyprogram Program
	jsonprogram, err := json.Marshal(program)
	if err != nil {
		log.Println("Could not convert program to JSON, is it in the correct format?")
	}
	rel := &url.URL{Path: "/api/programs"}
	u := c.BaseURL.ResolveReference(rel)
	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(jsonprogram))
	if err != nil {
		return emptyprogram, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return emptyprogram, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&program)
	return program, err
}

// UpdateProgram will update the specified program
func (c *Client) UpdateProgram(program Program, id string) (Program, error) {
	var emptyprogram Program
	jsonprogram, err := json.Marshal(program)
	if err != nil {
		log.Println("Could not convert program to JSON, is it in the correct format?")
	}
	rel := &url.URL{Path: "/api/programs/" + id}
	u := c.BaseURL.ResolveReference(rel)
	req, err := http.NewRequest("PUT", u.String(), bytes.NewBuffer(jsonprogram))
	if err != nil {
		return emptyprogram, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return emptyprogram, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&program)
	return program, err
}

// DeleteProgram will get a program
func (c *Client) DeleteProgram(id string) (bool, error) {
	rel := &url.URL{Path: "/api/programs/" + id}
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
	var program Program
	err = json.NewDecoder(resp.Body).Decode(&program)
	return true, err
}

// GetAssociatedRootDomains will get root domains belonging to the specified program
func (c *Client) GetAssociatedRootDomains(id string) ([]RootDomain, error) {
	rel := &url.URL{Path: "/api/programs/" + id + "/rootdomains"}
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

// GetAssociatedIPs will get IP addresses belonging to the specified program
func (c *Client) GetAssociatedIPs(id string) ([]IP, error) {
	rel := &url.URL{Path: "/api/programs/" + id + "/ips"}
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

// GetAssociatedVulns will get Vulns belonging to the specified program
func (c *Client) GetAssociatedVulns(id string) ([]Vuln, error) {
	rel := &url.URL{Path: "/api/programs/" + id + "/vulns"}
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

// GetAssociatedSubdomainsProgram will get Subdomains belonging to the specified program
func (c *Client) GetAssociatedSubdomainsProgram(id string) ([]Subdomain, error) {
	rel := &url.URL{Path: "/api/programs/" + id + "/subdomains"}
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
	var ips []Subdomain
	err = json.NewDecoder(resp.Body).Decode(&ips)
	return ips, err
}

// PrintPrograms prints the program to terminal in desired output format
func PrintPrograms(outputFormat string, c Client) {
	// get the programs
	programs, err := c.GetPrograms()
	if err != nil {
		fmt.Println("Error occured while fetching programs.", err)
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

// PrintProgram prints the program to terminal in desired output format
func PrintProgram(programID string, outputFormat string, c Client) {
	// get the program
	program, err := c.GetProgram(programID)
	if err != nil {
		fmt.Println("Error occured while fetching program.", err)
	}

	// output in desired format
	if outputFormat == "json" {

		programJSON, err := json.Marshal(program)
		if err != nil {
			fmt.Println("Error occured while converting the response to JSON: ", err)
		}
		fmt.Println(string(programJSON))

	} else {

		fmt.Println(program.ID)

	}
}

// PrintAssociatedRootDomains will print all rootdomains associated with the program
func PrintAssociatedRootDomains(id string, outputFormat string, c Client) {
	rootdomains, err := c.GetAssociatedRootDomains(id)
	if err != nil {
		fmt.Println("Error occured while fetching associated rootdomains: ", err)
		return
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

// PrintAssociatedSubdomainsProgram will print all subdomains associated with the program
func PrintAssociatedSubdomainsProgram(id string, outputFormat string, c Client) {
	subdomains, err := c.GetAssociatedSubdomainsProgram(id)
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

// PrintAssociatedIPs will print all ips associated with the program
func PrintAssociatedIPs(id string, outputFormat string, c Client) {
	ips, err := c.GetAssociatedIPs(id)
	if err != nil {
		fmt.Println("Error occured while fetching associated ips: ", err)
		return
	}
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

// PrintAssociatedVulns will print all ips associated with the program
func PrintAssociatedVulns(id string, outputFormat string, c Client) {
	vulns, err := c.GetAssociatedVulns(id)
	if err != nil {
		fmt.Println("Error occured while fetching associated vulns: ", err)
		return
	}
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

// ProgramCLI handles the program subcommand CLI
func ProgramCLI(c Client) {
	programsFlagSet := flag.NewFlagSet("programs", flag.ExitOnError)
	if len(os.Args) < 3 {
		fmt.Println("Invalid arguments. Hint: ./hakstore-client programs {list|create|delete}")
		return
	}
	switch os.Args[2] {
	case "list":
		programID := programsFlagSet.String("id", "", "ID of program")
		outputFormat := programsFlagSet.String("output", "", "output format")
		platformID := programsFlagSet.String("platform", "", "ID of platform")
		programsFlagSet.Parse(os.Args[3:])
		if isFlagPassed("id", programsFlagSet) {
			// show single program
			PrintProgram(*programID, *outputFormat, c)

		} else if isFlagPassed("platform", programsFlagSet) {
			PrintAssociatedPrograms(*platformID, *outputFormat, c)
		} else {
			// list programs
			PrintPrograms(*outputFormat, c)
		}
	case "create":
		programID := programsFlagSet.String("id", "", "ID of program")
		platformID := programsFlagSet.String("platform", "", "Platform that program is associated with")
		programsFlagSet.Parse(os.Args[3:])
		// create/update program
		if *programID == "" || *platformID == "" {
			fmt.Println("You need to specify a -id and -platform to create the program.")
			return
		}
		program := Program{ID: *programID, PlatformID: *platformID}
		_, err := c.CreateProgram(program)
		if err != nil {
			fmt.Println("An error occured while creating the program: ", err)
		}
	case "delete":
		programID := programsFlagSet.String("id", "", "ID of program")
		programsFlagSet.Parse(os.Args[3:])
		// delete program
		if *programID == "" {
			fmt.Println("You need to specify a program id to delete with -id.")
			return
		}
		c.DeleteProgram(*programID)

	// no valid subcommand found - default to showing a message and exiting
	default:
		fmt.Println("Invalid subsubcommand, ./hakstore-client programs {list|create|delete}")
		os.Exit(1)
	}
}
