package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// RootDomain is a structure to store details about bug bounty rootdomains
type RootDomain struct {
	gorm.Model
	ID         string      `json:"id" gorm:"PrimaryKey"`
	ProgramID  string      `json:"program"`
	Subdomains []Subdomain `json:"subdomains"`
}

// Get all RootDomains
func getRootDomains(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var rootdomains []RootDomain
	db.Find(&rootdomains)
	json.NewEncoder(w).Encode(rootdomains)
}

// Get a rootdomain
func getRootDomain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var rootdomain RootDomain
	vars := mux.Vars(r)
	db.Where("ID = ?", vars["id"]).Preload("Subdomains").Find(&rootdomain)
	json.NewEncoder(w).Encode(&rootdomain)
}

// Create a new rootdomain
func createRootDomain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var rootdomain RootDomain
	_ = json.NewDecoder(r.Body).Decode(&rootdomain)
	db.FirstOrCreate(&rootdomain)
	json.NewEncoder(w).Encode(rootdomain)
}

// Updates a rootdomain
func updateRootDomain(w http.ResponseWriter, r *http.Request) {
	//TODO
}

// Deletes a rootdomain
func deleteRootDomain(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var rootdomain RootDomain
	db.Where("ID = ?", id).Find(&rootdomain)
	deleteRootDomainLocal(rootdomain)
	var rootdomains []RootDomain
	db.Find(&rootdomains)
	json.NewEncoder(w).Encode(rootdomains)
}

func deleteRootDomainLocal(rootdomain RootDomain) {
	var subdomains []Subdomain
	db.Model(&rootdomain).Association("Subdomains").Find(&subdomains)
	for _, s := range subdomains {
		deleteSubdomainLocal(s)
	}
	db.Unscoped().Delete(&rootdomain)
}

// Dumps all subdomains associated with this rootdomain
func getAssociatedSubdomains(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var rootdomain RootDomain
	var subdomains []Subdomain
	vars := mux.Vars(r)
	db.Where("ID = ?", vars["id"]).Find(&rootdomain)
	db.Model(&rootdomain).Association("Subdomains").Find(&subdomains)
	json.NewEncoder(w).Encode(&subdomains)
}
