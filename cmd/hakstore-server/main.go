package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/hakluke/tldomains"
	"github.com/kelseyhightower/envconfig"
	"github.com/slack-go/slack"
	"gopkg.in/yaml.v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Config struct links with config.yml and environment variables
type Config struct {
	Server struct {
		Host string `yaml:"host" envconfig:"SERVER_HOST"`
		Port string `yaml:"port" envconfig:"SERVER_PORT"`
	} `yaml:"server"`
	Database struct {
		Host     string `yaml:"host" envconfig:"DATABASE_HOST"`
		User     string `yaml:"user" envconfig:"DATABASE_USER"`
		Password string `yaml:"password" envconfig:"DATABASE_PASSWORD"`
		DBName   string `yaml:"dbname" envconfig:"DATABASE_NAME"`
		Port     string `yaml:"port" envconfig:"DATABASE_PORT"`
		SSLmode  string `yaml:"sslmode" envconfig:"SSLMODE"`
		TimeZone string `yaml:"timezone" envconfig:"TIMEZONE"`
	} `yaml:"database"`
	Redis struct {
		Host     string `yaml:"host" envconfig:"REDIS_HOST"`
		Port     string `yaml:"port" envconfig:"REDIS_PORT"`
		Password string `yaml:"password" envconfig:"REDIS_PASSWORD"`
	} `yaml:"redis"`
	Slack struct {
		Webhook              string `yaml:"webhook" envconfig:"SLACK_WEBHOOK"`
		CriticalWebhook      string `yaml:"criticalwebhook" envconfig:"CRITICAL_WEBHOOK"`
		HighWebhook          string `yaml:"highwebhook" envconfig:"HIGH_WEBHOOK"`
		MediumWebhook        string `yaml:"mediumwebhook" envconfig:"MEDIUM_WEBHOOK"`
		LowWebhook           string `yaml:"lowwebhook" envconfig:"LOW_WEBHOOK"`
		InformationalWebhook string `yaml:"informationalwebhook" envconfig:"INFORMATIONAL_WEBHOOK"`
	} `yaml:"slack"`
}

// GLOBAL VARIABLES

// Database
var db *gorm.DB

// Redis client
var redisClient *redis.Client

// Cache file for TLDs
var tldomainsCache *tldomains.Tldomains

// Config struct
var config *Config

// Slack client
var slackAPI *slack.Client

func main() {
	fmt.Println("Starting hakstore...")

	var err error

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

	// Start mux and define the routes in routes.go
	r := mux.NewRouter()
	defineRoutes(r)

	// Create cache file of TLDs
	tldomainsCache, err = tldomains.New("/tmp/tld.cache")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating cache file for TLDs: %s", err)
		os.Exit(1)
	}

	// Connect to redis server
	redisClient = redis.NewClient(&redis.Options{
		Addr:     config.Redis.Host + ":" + config.Redis.Port,
		Password: config.Redis.Password,
		DB:       0, // redis databases are deprecated so we will just use the default
	})

	// Connect to the database
	dsn := "host=" + config.Database.Host + " user=" + config.Database.User + " password=" + config.Database.Password + " dbname=" + config.Database.DBName + " port=" + config.Database.Port + " sslmode=" + config.Database.SSLmode + " timezone=" + config.Database.TimeZone
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database")
	}

	// Start authentication middleware
	keyUsers := make(map[string]string) // initialise an empty map
	amw := authenticationMiddleware{}
	amw.keyUsers = keyUsers
	amw.Populate() // populate keyusers map with actual api keys from DB
	r.Use(amw.Middleware)

	// Migrate the schema
	db.AutoMigrate(&Platform{}, &Program{}, &RootDomain{}, &Subdomain{}, &IP{}, &User{}, &Vuln{})

	// If no users exist yet, create the first one!
	var user User
	user.ID = "admin"
	db.FirstOrCreate(&user)
	fmt.Println("Your admin API Key is", user.Key)
	fmt.Println("Keep it secret, keep it safe.")

	// Populate the database with some test data
	//seedDB()

	var subdomains []Subdomain
	db.Find(&subdomains)

	// Set up the flags
	portPtr := flag.Uint("serve", 80, "starts hakstore server on the specified port")
	flag.Parse()

	// Start the web server
	fmt.Println("Starting web server on port " + fmt.Sprint(*portPtr))
	log.Fatal(http.ListenAndServe(":"+fmt.Sprint(*portPtr), r))

}

// Reads environment variables into the config struct (for docker)
func readEnv(cfg *Config) {
	err := envconfig.Process("", cfg)
	if err != nil {
		fmt.Println("You need to either set up config.yml or add environment variables.")
		os.Exit(1)
	}
}

// Populates the DB with some initial data (mostly for testing purposes)
func seedDB() {
	fmt.Println("Seeding the database with some dummy test data...")
	// Create some dummy data
	var platforms = []Platform{
		{ID: "bugcrowd", URL: "https://bugcrowd.com"},
		{ID: "hackerone", URL: "https://hackerone.com"},
		{ID: "intigriti", URL: "https://intigriti.com"},
	}

	var programs = []Program{
		{ID: "tesla", PlatformID: "bugcrowd"},
		{ID: "tripadvisor", PlatformID: "bugcrowd"},
	}

	var rootdomains = []RootDomain{
		{ID: "tesla.com", ProgramID: "tesla"},
		{ID: "tripadvisor.com", ProgramID: "tripadvisor"},
	}

	var subdomains = []Subdomain{
		{ID: "www.tesla.com", RootDomainID: "tesla.com"},
		{ID: "www.tripadvisor.com", RootDomainID: "tripadvisor.com"},
	}

	// Push the dummy data to the database
	db.Create(&platforms)
	db.Create(&programs)
	db.Create(&rootdomains)
	db.Create(&subdomains)
}
