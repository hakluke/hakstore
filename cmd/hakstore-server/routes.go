package main

import (
	"github.com/gorilla/mux"
)

func defineRoutes(r *mux.Router) {
	// Platform routes
	r.HandleFunc("/api/platforms", getPlatforms).Methods("GET")
	r.HandleFunc("/api/platforms", createPlatform).Methods("POST")
	r.HandleFunc("/api/platforms/{id}", getPlatform).Methods("GET")
	r.HandleFunc("/api/platforms/{id}", updatePlatform).Methods("PUT")
	r.HandleFunc("/api/platforms/{id}", deletePlatform).Methods("DELETE")
	r.HandleFunc("/api/platforms/{id}/programs", getAssociatedPrograms).Methods("GET")

	// Program routes
	r.HandleFunc("/api/programs", getPrograms).Methods("GET")
	r.HandleFunc("/api/programs", createProgram).Methods("POST")
	r.HandleFunc("/api/programs/{id}", getProgram).Methods("GET")
	r.HandleFunc("/api/programs/{id}", updateProgram).Methods("PUT")
	r.HandleFunc("/api/programs/{id}", deleteProgram).Methods("DELETE")
	r.HandleFunc("/api/programs/{id}/rootdomains", getAssociatedRootDomains).Methods("GET")
	r.HandleFunc("/api/programs/{id}/ips", getAssociatedIPs).Methods("GET")
	r.HandleFunc("/api/programs/{id}/subdomains", getAssociatedSubdomainsProgram).Methods("GET")

	// RootDomain routes
	r.HandleFunc("/api/rootdomains", getRootDomains).Methods("GET")
	r.HandleFunc("/api/rootdomains", createRootDomain).Methods("POST")
	r.HandleFunc("/api/rootdomains/{id}", getRootDomain).Methods("GET")
	r.HandleFunc("/api/rootdomains/{id}", updateRootDomain).Methods("PUT")
	r.HandleFunc("/api/rootdomains/{id}", deleteRootDomain).Methods("DELETE")
	r.HandleFunc("/api/rootdomains/{id}/subdomains", getAssociatedSubdomains).Methods("GET")

	// Subdomain routes
	r.HandleFunc("/api/subdomains", getSubdomains).Methods("GET")
	r.HandleFunc("/api/subdomains", createSubdomains).Methods("POST")
	r.HandleFunc("/api/subdomains/{id}", getSubdomain).Methods("GET")
	r.HandleFunc("/api/subdomains/{id}", updateSubdomain).Methods("PUT")
	r.HandleFunc("/api/subdomains/{id}", deleteSubdomain).Methods("DELETE")
	r.HandleFunc("/api/subdomains/recent/{minutes}", getRecentSubdomains).Methods("GET")
	r.HandleFunc("/api/subdomains/{id}/ips", associateIPWithSubdomain).Methods("POST")

	// IP routes
	r.HandleFunc("/api/ips", getIPs).Methods("GET")
	r.HandleFunc("/api/ips", createIPs).Methods("POST")
	r.HandleFunc("/api/ips/{id}", getIP).Methods("GET")
	r.HandleFunc("/api/ips/{id}", updateIP).Methods("PUT")
	r.HandleFunc("/api/ips/{id}", deleteIP).Methods("DELETE")

	// Vuln routes
	r.HandleFunc("/api/vulns", getVulns).Methods("GET")
	r.HandleFunc("/api/vulns", createVulns).Methods("POST")
	r.HandleFunc("/api/vulns/{id}", getVuln).Methods("GET")
	r.HandleFunc("/api/vulns/{id}", updateVuln).Methods("PUT")
	r.HandleFunc("/api/vulns/{id}", deleteVuln).Methods("DELETE")

	// Job routes
	r.HandleFunc("/api/jobs", createJobs).Methods("POST")
}
