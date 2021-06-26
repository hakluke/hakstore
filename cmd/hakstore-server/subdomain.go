package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Subdomain is a structure to store details about bug bounty subdomains
type Subdomain struct {
	gorm.Model
	ID           string `json:"id" gorm:"PrimaryKey"`
	ProgramID    string `json:"program"`
	RootDomainID string `json:"rootdomain"`
	CNAME        string `json:"cname"`
	Nameservers  string `json:"nameservers"`
	IPs          []*IP  `json:"ips" gorm:"many2many:subdomain_ips;"`
}

// BeforeCreate will associate the subdomain to the appropriate rootdomain, unless a specific rootdomain is specified
func (s *Subdomain) BeforeCreate(tx *gorm.DB) (err error) {
	// try to automatically determine the rootdomain if none are supplied
	// if s.RootDomainID == "" {
	// 	extract := tldomainsCache.Parse(s.ID)
	// 	s.RootDomainID = extract.Root + "." + extract.Suffix
	// 	var rootdomain RootDomain
	// 	db.Where("ID = ?", s.RootDomainID).FirstOrInit(&rootdomain)
	// 	rootdomain.ProgramID = s.ProgramID
	// 	db.Save(&rootdomain)
	// }

	// set the ProgramID to the same as the ProgramID on the associated rootdomain
	var rootdomain RootDomain
	db.Where("ID = ?", s.RootDomainID).FirstOrInit(&rootdomain)
	s.ProgramID = rootdomain.ProgramID
	return nil
}

// // AfterCreate will associate the subdomain to the appropriate rootdomain, unless a specific rootdomain is specified
// func (s *Subdomain) AfterCreate(tx *gorm.DB) (err error) {
// 	createJob("updateDNSData", s.ID)
// 	createJob("nuclei", s.ID)
// 	return nil
// }

// Get all Subdomains
func getSubdomains(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var subdomains []Subdomain
	db.Find(&subdomains)
	json.NewEncoder(w).Encode(subdomains)
}

// Get a specific subdomain
func getSubdomain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var subdomain Subdomain
	vars := mux.Vars(r)
	db.Preload("IPs").Where("ID = ?", vars["id"]).Find(&subdomain)
	json.NewEncoder(w).Encode(&subdomain)
}

// Creates new subdomains, accepts batches
func createSubdomains(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var subdomains []Subdomain
	_ = json.NewDecoder(r.Body).Decode(&subdomains)
	db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"root_domain_id"}),
	}).Create(&subdomains)
	json.NewEncoder(w).Encode(subdomains)
}

// Updates a subdomain
func updateSubdomain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var subdomain Subdomain
	_ = json.NewDecoder(r.Body).Decode(&subdomain)

	db.Model(&subdomain).Updates(map[string]interface{}{
		"root_domain_id": subdomain.RootDomainID,
		"nameservers":    subdomain.Nameservers,
		"cname":          subdomain.CNAME,
	})

	json.NewEncoder(w).Encode(subdomain)
}

// Deletes a subdomain
func deleteSubdomain(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var subdomain Subdomain
	db.Where("ID = ?", id).Find(&subdomain)
	deleteSubdomainLocal(subdomain)
	var subdomains []Subdomain
	db.Find(&subdomains)
	json.NewEncoder(w).Encode(subdomains)
}

func deleteSubdomainLocal(subdomain Subdomain) {
	db.Model(&subdomain).Association("IPs").Clear()
	db.Unscoped().Delete(&subdomain)
}

func getRecentSubdomains(w http.ResponseWriter, r *http.Request) {
	var subdomains []Subdomain
	var recentSubdomains []Subdomain
	vars := mux.Vars(r)
	minutes := vars["minutes"]
	minutesInt, err := strconv.Atoi(minutes)
	if err != nil {
		return
	}
	duration := time.Duration(minutesInt) * time.Minute
	db.Find(&subdomains)
	now := time.Now()
	for _, subdomain := range subdomains {
		sinceInception := now.Sub(subdomain.CreatedAt)
		if duration > sinceInception {
			recentSubdomains = append(recentSubdomains, subdomain)
		}
	}
	json.NewEncoder(w).Encode(recentSubdomains)
}

func getRecentSubdomainsLocal(minutes uint32) []Subdomain {
	var subdomains []Subdomain
	var recentSubdomains []Subdomain
	duration := time.Duration(minutes) * time.Minute
	db.Find(&subdomains)
	now := time.Now()
	for _, subdomain := range subdomains {
		sinceInception := now.Sub(subdomain.CreatedAt)
		if duration > sinceInception {
			recentSubdomains = append(recentSubdomains, subdomain)
		}
	}
	return recentSubdomains
}

func printRecentSubdomains(minutes uint32) string {
	subdomains := getRecentSubdomainsLocal(minutes)
	text := "New subdomains discovered in the last " + fmt.Sprint(minutes) + " minutes:\n```"
	for _, s := range subdomains {
		text = text + "- " + s.ID + "\n"
	}
	text = text + "```"
	fmt.Println(text)
	return text
}

func associateIPWithSubdomain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get the ips from the request body
	var ips []IP
	_ = json.NewDecoder(r.Body).Decode(&ips)

	// Get the subdomain from the URL
	var subdomain Subdomain
	vars := mux.Vars(r)
	db.Where("ID = ?", vars["id"]).Find(&subdomain)

	// For each IP, associate it with the given subdomain
	db.Model(&subdomain).Association("IPs")
	for _, ip := range ips {
		var rootdomain RootDomain
		db.Where("id = ?", subdomain.RootDomainID).First(&rootdomain)
		db.Where(IP{ID: ip.ID, ProgramID: rootdomain.ProgramID}).FirstOrCreate(&ip)
		db.Model(&subdomain).Association("IPs").Append([]IP{ip})
	}
	json.NewEncoder(w).Encode(ips)
}

// updateDNSData gets the IP addresses, nameservers and CNAME data from a specified subdomain and saves it.
func updateDNSData(subdomain Subdomain) (err error) {
	didReturnIP := true
	iprecords, tempErr := net.LookupIP(subdomain.ID)
	if tempErr != nil {
		didReturnIP = false
	}

	if didReturnIP {
		for _, ip := range iprecords {
			var newIP IP
			ipString := fmt.Sprint(ip)
			db.Model(&subdomain).Association("IPs")

			var rootdomain RootDomain
			db.Where("id = ?", subdomain.RootDomainID).First(&rootdomain)
			db.Where(IP{ID: ipString, ProgramID: rootdomain.ProgramID}).FirstOrCreate(&newIP)
			db.Model(&subdomain).Association("IPs").Append([]IP{newIP})
		}
	}

	// Update the CNAME if there is one
	didReturnCNAME := true
	cnamerecord, tempErr := net.LookupCNAME(subdomain.ID)
	if tempErr != nil {
		didReturnCNAME = false
	}
	if didReturnCNAME {
		if cnamerecord != subdomain.ID+"." {
			db.Model(&subdomain).Update("CNAME", cnamerecord)
			fmt.Println("added CNAME", cnamerecord, "to subdomain", subdomain.ID)
		}
	}

	didReturnNS := true
	nameservers, tempErr := net.LookupNS(subdomain.ID)
	if tempErr != nil {
		didReturnNS = false
	}
	if didReturnNS {
		var nameserverStrings []string
		for _, ns := range nameservers {
			nameserverStrings = append(nameserverStrings, ns.Host)
		}

		// Just keep the nameservers in a json string in a single text field in the DB... TODO add a new Model with relationships instead
		jsonNameservers, tempErr := json.Marshal(nameservers)
		if tempErr != nil {
			fmt.Println("Error marshalling nameservers:", tempErr)
		}
		db.Model(&subdomain).Update("Nameservers", string(jsonNameservers))
	}
	return nil
}
