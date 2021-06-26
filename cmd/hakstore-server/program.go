package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// Program is a structure to store details about bug bounty programs
type Program struct {
	gorm.Model
	ID          string       `json:"id" gorm:"PrimaryKey"`
	PlatformID  string       `json:"platform"`
	Subdomains  []Subdomain  `json:"subdomains"`
	RootDomains []RootDomain `json:"rootdomains"`
	IPs         []IP         `json:"ips"`
}

// Get all Programs
func getPrograms(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var programs []Program
	db.Find(&programs)
	json.NewEncoder(w).Encode(programs)
}

// Get a program
func getProgram(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var program Program
	vars := mux.Vars(r)
	db.Where("ID = ?", vars["id"]).Preload("RootDomains").Find(&program)
	json.NewEncoder(w).Encode(&program)
}

// Create a new program
func createProgram(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var program Program
	_ = json.NewDecoder(r.Body).Decode(&program)
	db.FirstOrCreate(&program)
	json.NewEncoder(w).Encode(program)
}

// Updates a program
func updateProgram(w http.ResponseWriter, r *http.Request) {
	//TODO
}

// Deletes a program
func deleteProgram(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var program Program
	db.Where("id = ?", id).Find(&program)
	deleteProgramLocal(program)
	var programs []Program
	db.Find(&programs)
	json.NewEncoder(w).Encode(programs)
}

func deleteProgramLocal(program Program) {
	db.Model(&program).Association("IPs").Clear()
	var rootdomains []RootDomain
	db.Model(&program).Association("RootDomains").Find(&rootdomains)
	for _, s := range rootdomains {
		deleteRootDomainLocal(s)
	}
	db.Unscoped().Delete(&program)
}

// Dumps all root domains associated with this platform
func getAssociatedRootDomains(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var program Program
	var rootdomains []RootDomain
	vars := mux.Vars(r)
	db.Where("ID = ?", vars["id"]).Find(&program)
	db.Model(&program).Association("RootDomains").Find(&rootdomains)
	json.NewEncoder(w).Encode(&rootdomains)
}

// Dumps all IPs associated with this program
func getAssociatedIPs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var program Program
	var ips []IP
	vars := mux.Vars(r)
	db.Where("ID = ?", vars["id"]).Find(&program)
	db.Model(&program).Association("IPs").Find(&ips)
	json.NewEncoder(w).Encode(&ips)
}

// Dumps all Subdomains associated with this program
func getAssociatedSubdomainsProgram(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var program Program
	var subdomains []Subdomain
	vars := mux.Vars(r)
	db.Where("ID = ?", vars["id"]).Find(&program)
	db.Model(&program).Association("Subdomains").Find(&subdomains)
	json.NewEncoder(w).Encode(&subdomains)
}
