package hakstoreclient

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

// Subdomain is a structure to store details about bug bounty subdomains
type Subdomain struct {
	gorm.Model
	ID           string `json:"id" gorm:"PrimaryKey"`
	RootDomainID string `json:"rootdomain"`
	CNAME        string `json:"cname"`
	Nameservers  string `json:"nameservers"`
	IPs          []*IP  `json:"ips" gorm:"many2many:subdomain_ips;"`
}

// GetSubdomains will get all subdomains from database
func (c *Client) GetSubdomains() ([]Subdomain, error) {
	rel := &url.URL{Path: "/api/subdomains"}
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

// GetSubdomain will get a subdomain
func (c *Client) GetSubdomain(id string) (Subdomain, error) {

	var emptysubdomain Subdomain
	rel := &url.URL{Path: "/api/subdomains/" + id}
	u := c.BaseURL.ResolveReference(rel)
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return emptysubdomain, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return emptysubdomain, err
	}
	defer resp.Body.Close()
	var subdomain Subdomain
	err = json.NewDecoder(resp.Body).Decode(&subdomain)
	return subdomain, err
}

// CreateSubdomains will create a new subdomain
func (c *Client) CreateSubdomains(subdomains []Subdomain) ([]Subdomain, error) {
	var emptysubdomain []Subdomain

	jsonsubdomains, err := json.Marshal(subdomains)
	if err != nil {
		log.Println("Could not convert subdomain to JSON, is it in the correct format?")
	}
	rel := &url.URL{Path: "/api/subdomains"}
	u := c.BaseURL.ResolveReference(rel)
	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(jsonsubdomains))
	if err != nil {
		return emptysubdomain, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return emptysubdomain, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&subdomains)
	return subdomains, err
}

// UpdateSubdomains updates all subdomains in a list. It simply calls UpdateSubdomain multiple times
func (c *Client) UpdateSubdomains(subdomains []Subdomain) {
	for _, subdomain := range subdomains {
		c.UpdateSubdomain(subdomain, subdomain.ID)
	}
}

// UpdateSubdomain will update the specified subdomain
func (c *Client) UpdateSubdomain(subdomain Subdomain, id string) (Subdomain, error) {
	var emptysubdomain Subdomain
	jsonsubdomain, err := json.Marshal(subdomain)
	if err != nil {
		log.Println("Could not convert subdomain to JSON, is it in the correct format?")
	}
	rel := &url.URL{Path: "/api/subdomains/" + id}
	u := c.BaseURL.ResolveReference(rel)
	req, err := http.NewRequest("PUT", u.String(), bytes.NewBuffer(jsonsubdomain))
	if err != nil {
		return emptysubdomain, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return emptysubdomain, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&subdomain)
	return subdomain, err
}

// DeleteSubdomain will delete a subdomain
func (c *Client) DeleteSubdomain(id string) (bool, error) {
	rel := &url.URL{Path: "/api/subdomains/" + id}
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

// GetRecentSubdomains will get subdomains that have been created within the last x minutes
func (c *Client) GetRecentSubdomains(minutes uint32) ([]Subdomain, error) {
	rel := &url.URL{Path: "/api/subdomains/recent/" + strconv.FormatUint(uint64(minutes), 10)}
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

// PrintSubdomains prints multiple subdomains to terminal in desired output format
func PrintSubdomains(outputFormat string, subdomains []Subdomain) {

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

// PrintSubdomain prints the subdomain to terminal in desired output format
func PrintSubdomain(subdomainID string, outputFormat string, c Client) {
	// get the subdomain
	subdomain, err := c.GetSubdomain(subdomainID)
	if err != nil {
		fmt.Println("Error occured while fetching subdomain.", err)
	}

	// output in desired format
	if outputFormat == "json" {

		subdomainJSON, err := json.Marshal(subdomain)
		if err != nil {
			fmt.Println("Error occured while converting the response to JSON: ", err)
		}
		fmt.Println(string(subdomainJSON))

	} else {
		fmt.Println(subdomain.ID)
	}
}

// ImportSubdomainsFromFile adds each line of a file full of subdomains and associates it with the provided rootdomain
func ImportSubdomainsFromFile(filename string, rootdomain string, c Client) {
	var subdomains []Subdomain
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("Error opening file: ", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var newSubdomain Subdomain
	for scanner.Scan() {
		newSubdomain = Subdomain{
			ID:           scanner.Text(),
			RootDomainID: rootdomain,
		}
		subdomains = append(subdomains, newSubdomain)
	}

	c.CreateSubdomains(subdomains)

	if err := scanner.Err(); err != nil {
		log.Fatal("Error scanning file:", err)
	}

}

// AssociateIPWithSubdomain will create new IPs and associate them with the given subdomain
func (c *Client) AssociateIPWithSubdomain(ips []IP, subdomain Subdomain) ([]IP, error) {
	var emptyip []IP

	jsonips, err := json.Marshal(ips)
	if err != nil {
		log.Println("Could not convert ip to JSON, is it in the correct format?", err)
	}
	rel := &url.URL{Path: "/api/subdomains/" + subdomain.ID + "/ips"}
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

// SubdomainsCLI handles the subdomains subcommand CLI
func SubdomainsCLI(c Client) {
	if len(os.Args) < 3 {
		fmt.Println("Invalid arguments. Hint: ./hakstore-client subdomains {list|create|delete|associateips|import}")
		return
	}
	switch os.Args[2] {
	case "list":
		subdomainsFlagSet := flag.NewFlagSet("subdomains list", flag.ExitOnError)
		subdomainID := subdomainsFlagSet.String("id", "", "ID of subdomain")
		outputFormat := subdomainsFlagSet.String("output", "", "output format")
		rootdomainID := subdomainsFlagSet.String("rootdomain", "", "ID of rootdomain")
		programID := subdomainsFlagSet.String("program", "", "ID of program")
		recent := subdomainsFlagSet.Int("recent", 0, "number of minutes")
		subdomainsFlagSet.Parse(os.Args[3:])
		if isFlagPassed("id", subdomainsFlagSet) {
			// show single subdomain
			PrintSubdomain(*subdomainID, *outputFormat, c)

		} else if isFlagPassed("rootdomain", subdomainsFlagSet) {
			PrintAssociatedSubdomains(*rootdomainID, *outputFormat, c)
		} else if isFlagPassed("program", subdomainsFlagSet) {
			PrintAssociatedSubdomainsProgram(*programID, *outputFormat, c)
		} else if isFlagPassed("recent", subdomainsFlagSet) {
			subdomains, err := c.GetRecentSubdomains(uint32(*recent))
			if err != nil {
				fmt.Println("Error encountered while getting recent subdomains:", err)
			}
			PrintSubdomains(*outputFormat, subdomains)
		} else {
			// list subdomains
			subdomains, err := c.GetSubdomains()
			if err != nil {
				fmt.Println("Error retreiving subdomains: ", err)
			}
			PrintSubdomains(*outputFormat, subdomains)
		}
	case "create":
		subdomainsFlagSet := flag.NewFlagSet("subdomains create", flag.ExitOnError)
		subdomainID := subdomainsFlagSet.String("id", "", "ID of subdomain, e.g. api.example.com")
		rootdomainID := subdomainsFlagSet.String("rootdomain", "", "Rootdomain that the subdomain will be associated with, e.g. example.com. If you don't specify a rootdomain, it will be extracted from the subdomain.")
		subdomainsFlagSet.Parse(os.Args[3:])
		// create/update subdomain
		if *subdomainID == "" {
			fmt.Println("You need to specify a -id and a -rootdomain to create the subdomain.")
			return
		}
		subdomains := []Subdomain{{ID: *subdomainID, RootDomainID: *rootdomainID}}
		_, err := c.CreateSubdomains(subdomains)
		if err != nil {
			fmt.Println("An error occured while creating the subdomain: ", err)
		}
	case "delete":
		subdomainsFlagSet := flag.NewFlagSet("subdomains delete", flag.ExitOnError)
		subdomainID := subdomainsFlagSet.String("id", "", "ID of subdomain")
		subdomainsFlagSet.Parse(os.Args[3:])
		// delete subdomain
		if *subdomainID == "" {
			fmt.Println("You need to specify a subdomain id to delete with -id.")
			return
		}
		c.DeleteSubdomain(*subdomainID)

	case "associateips":
		subdomainsFlagSet := flag.NewFlagSet("associate ips", flag.ExitOnError)
		subdomainID := subdomainsFlagSet.String("subdomain", "", "ID of subdomain")
		ipString := subdomainsFlagSet.String("ips", "", "comma separated list of IP addresses")
		subdomainsFlagSet.Parse(os.Args[3:])
		var ips []IP
		subdomain, err := c.GetSubdomain(*subdomainID)
		if err != nil {
			log.Fatal("No subdomain exists with the specified ID.", err)
		}
		rootdomain, err := c.GetRootDomain(subdomain.RootDomainID)
		if err != nil {
			log.Fatal("Not able to fetch root domain belonging to the specified subdomain.", err)
		}
		// create IP objects out of the comma separated string
		ipStrings := strings.Split(*ipString, ",")
		for _, ip := range ipStrings {
			ips = append(ips, IP{ID: ip, ProgramID: rootdomain.ProgramID})
		}
		c.AssociateIPWithSubdomain(ips, subdomain)

	case "import":
		subdomainsFlagSet := flag.NewFlagSet("platforms", flag.ExitOnError)
		file := subdomainsFlagSet.String("file", "", "File that you wish to import from")
		rootdomain := subdomainsFlagSet.String("rootdomain", "", "Root domain of the subdomains to import")
		subdomainsFlagSet.Parse(os.Args[3:])
		if isFlagPassed("file", subdomainsFlagSet) && isFlagPassed("rootdomain", subdomainsFlagSet) {
			ImportSubdomainsFromFile(*file, *rootdomain, c)
		} else {
			fmt.Println("Usage: hakstore-client subdomains import -file <./subs.txt> -rootdomain example.com")
		}

	// no valid subcommand found - default to showing a message and exiting
	default:
		fmt.Println("Invalid subsubcommand, ./hakstore-client subdomains {list|create|delete}")
		os.Exit(1)
	}
}
