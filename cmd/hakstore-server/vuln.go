package main

import (
	"encoding/json"
	"net/http"

	"gorm.io/gorm"

	"github.com/gorilla/mux"
)

// Vuln is a structure to store details about vulnerabilities
type Vuln struct {
	gorm.Model
	ID          int          `json:"id" gorm:"PrimaryKey;autoIncrement"`
	Subdomains  []*Subdomain `json:"subdomains" gorm:"many2many:subdomain_vulns;"`
	IPs         []*IP        `json:"ips" gorm:"many2many:ip_vulns;"`
	Description string       `json:"description"`
	ProgramID   string       `json:"program"`
	Severity    int          `json:"severity"`
}

func severityString(severity int) string {
	switch severity {
	case 1:
		return "critical"
	case 2:
		return "high"
	case 3:
		return "medium"
	case 4:
		return "low"
	case 5:
		return "informational"
	}
	return "unknown" // we shouldn't ever have a Vuln without one of the severities in the switch, so we shouldn't get here.
}

// BeforeCreate will add the programID if it is missing
func (v *Vuln) BeforeCreate(tx *gorm.DB) (err error) {
	// fill the program field based on the subdomain
	if v.ProgramID == "" {
		var sub Subdomain
		db.Where("ID = ?", v.Subdomains[0].ID).FirstOrInit(&sub)
		v.ProgramID = sub.ProgramID
	}
	return nil
}

// AfterCreate sends a Slack message with the vuln details
func (v *Vuln) AfterCreate(tx *gorm.DB) (err error) {
	// send a notification about the vuln
	SendVulnNotification(*v)
	return nil
}

// Get all Vulns
func getVulns(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var vulns []Vuln
	db.Find(&vulns)
	json.NewEncoder(w).Encode(vulns)
}

// Get a specific vuln
func getVuln(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var vuln Vuln
	vars := mux.Vars(r)
	db.Where("ID = ?", vars["id"]).Find(&vuln)
	json.NewEncoder(w).Encode(&vuln)
}

// Creates new vulns, accepts batches
func createVulns(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var vulns []Vuln
	_ = json.NewDecoder(r.Body).Decode(&vulns)
	db.Create(&vulns)
	json.NewEncoder(w).Encode(vulns)
}

// Updates a vuln
func updateVuln(w http.ResponseWriter, r *http.Request) {
	// not necessary, let's just delete it and create a new one
}

// Deletes a vuln
func deleteVuln(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var vuln Vuln
	db.Where("ID = ?", id).Find(&vuln)
	deleteVulnLocal(vuln)
	var vulns []Vuln
	db.Find(&vulns)
	json.NewEncoder(w).Encode(vulns)
}

// Deletes relationships and then removes the model
func deleteVulnLocal(vuln Vuln) {
	db.Model(&vuln).Association("Subdomains").Clear()
	db.Model(&vuln).Association("IPs").Clear()
	db.Unscoped().Delete(&vuln)
}
