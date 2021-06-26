package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/hakluke/haktools/pkg/hakstoreclient"
	"gopkg.in/yaml.v2"
)

var c hakstoreclient.Client

// Implement a "round tripper" to add an API key to every request
type transport struct {
	underlyingTransport http.RoundTripper
}

// The function that is performed for the round tripper (adds an api key header)
func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("X-API-Key", config.Server.APIKey)
	return t.underlyingTransport.RoundTrip(req)
}

// Config struct
type Config struct {
	Server struct {
		BaseURL string `yaml:"baseurl" envconfig:"HAKSTORE_BASEURL"`
		APIKey  string `yaml:"apikey" envconfig:"API_KEY"`
	} `yaml:"server"`
}

var config *Config

func main() {

	// load config file
	f, err := os.Open(os.Getenv("HOME") + "/.config/haktools/hakstore-config.yml")
	if err != nil {
		fmt.Println("Error opening config file, does it exist?:", err)
	}
	defer f.Close()

	// parse the config file
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("Error decoding config.yml", err)
		return
	}

	url, err := url.Parse(config.Server.BaseURL)
	if err != nil {
		panic("Base URL provided is not valid.")
	}

	// Use the roundtripper
	httpclient := http.Client{Transport: &transport{underlyingTransport: http.DefaultTransport}}

	c := hakstoreclient.Client{
		BaseURL:    url,
		UserAgent:  "hakstore client", //TODO get this from a config file
		HTTPClient: &httpclient,
	}

	if len(os.Args) < 2 {
		fmt.Println("Subcommand missing or incorrect. Hint: hakstore-client {platforms|programs|rootdomains|subdomains|import}")
		os.Exit(1)
	}

	// parse the cli arguments

	switch os.Args[1] {
	case "export":
		hakstoreclient.ExportCLI(c)
	case "platforms":
		hakstoreclient.PlatformCLI(c)
	case "programs":
		hakstoreclient.ProgramCLI(c)
	case "rootdomains":
		hakstoreclient.RootdomainsCLI(c)
	case "subdomains":
		hakstoreclient.SubdomainsCLI(c)
	case "ips":
		hakstoreclient.IPsCLI(c)
	case "vulns":
		hakstoreclient.VulnsCLI(c)
	case "jobs":
		hakstoreclient.JobsCLI(c)
	// no valid subcommand found - default to showing a message and exiting
	default:
		fmt.Println("Subcommand missing or incorrect. Hint: hakstore-client {platforms|programs|rootdomains|subdomains|ips|vulns|jobs}")
		os.Exit(1)
	}
}

// isFlagPassed checks whether a cli flag has been used
func isFlagPassed(name string, flagset *flag.FlagSet) bool {
	found := false
	flagset.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
