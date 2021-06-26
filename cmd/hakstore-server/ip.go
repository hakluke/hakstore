package main

import (
	"encoding/json"
	"net/http"

	"gorm.io/gorm"

	"github.com/gorilla/mux"
)

// IP is a structure to store details about ip addresses
type IP struct {
	gorm.Model
	ID         string       `json:"id" gorm:"PrimaryKey"`
	Subdomains []*Subdomain `json:"subdomains" gorm:"many2many:subdomain_ips;"`
	ProgramID  string       `json:"program"`
}

// Get all IPs
func getIPs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var ips []IP
	db.Find(&ips)
	json.NewEncoder(w).Encode(ips)
}

// Get a specific ip
func getIP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var ip IP
	vars := mux.Vars(r)
	db.Where("ID = ?", vars["id"]).Find(&ip)
	json.NewEncoder(w).Encode(&ip)
}

// Creates new ips, accepts batches
func createIPs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var ips []IP
	_ = json.NewDecoder(r.Body).Decode(&ips)
	db.Create(&ips)
	json.NewEncoder(w).Encode(ips)
}

// Updates a ip
func updateIP(w http.ResponseWriter, r *http.Request) {
	// not necessary because createIP() uses UPSERT
}

// Deletes a ip
func deleteIP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var ip IP
	db.Where("ID = ?", id).Find(&ip)
	deleteIPLocal(ip)
	var ips []IP
	db.Find(&ips)
	json.NewEncoder(w).Encode(ips)
}

func deleteIPLocal(ip IP) {
	db.Model(&ip).Association("Subdomains").Clear()
	db.Unscoped().Delete(&ip)
}
