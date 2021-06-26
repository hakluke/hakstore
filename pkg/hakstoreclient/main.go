package hakstoreclient

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

// Client for interacting with hakstore server
type Client struct {
	BaseURL    *url.URL
	UserAgent  string
	HTTPClient *http.Client
}

func main() {
	url, err := url.Parse("http://localhost") //TODO get this from a config file
	if err != nil {
		panic("Base URL provided is not valid.")
	}
	httpclient := http.Client{}
	c := Client{
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
		ExportCLI(c)
	case "platforms":
		PlatformCLI(c)
	case "programs":
		ProgramCLI(c)
	case "rootdomains":
		RootdomainsCLI(c)
	case "subdomains":
		SubdomainsCLI(c)
	case "ips":
		IPsCLI(c)
	case "jobs":
		JobsCLI(c)
	// no valid subcommand found - default to showing a message and exiting
	default:
		fmt.Println("Subcommand missing or incorrect. Hint: hakstore-client {platforms|programs|rootdomains|subdomains}")
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
