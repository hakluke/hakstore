package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Platform Struct (Model)
type Platform struct {
	gorm.Model
	ID       string    `json:"id" gorm:"PrimaryKey"`
	URL      string    `json:"url"`
	Programs []Program `json:"programs"`
}

// Get all platforms
func getPlatforms(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var platforms []Platform
	db.Find(&platforms)
	json.NewEncoder(w).Encode(platforms)
}

// Get a platform
func getPlatform(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var platform Platform
	vars := mux.Vars(r)
	db.Where("ID = ?", vars["id"]).Preload("Programs").Find(&platform)
	json.NewEncoder(w).Encode(&platform)
}

// Create a new platform
func createPlatform(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var platform Platform
	var blankDeletedAt gorm.DeletedAt
	_ = json.NewDecoder(r.Body).Decode(&platform)
	platform.DeletedAt = blankDeletedAt
	db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"url", "deleted_at"}),
	}).FirstOrCreate(&platform)
	json.NewEncoder(w).Encode(platform)
}

// Updates a platform
func updatePlatform(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	url := vars["url"]
	var platform Platform
	db.Where("id = ?", id).Find(&platform)
	platform.URL = url
	db.Save(&platform)
}

// Deletes a platform
func deletePlatform(w http.ResponseWriter, r *http.Request) {
	// get platform from request
	vars := mux.Vars(r)
	id := vars["id"]
	var platform Platform
	db.Where("id = ?", id).Find(&platform)
	deletePlatformLocal(platform)
	var platforms []Platform
	db.Find(&platforms)
	json.NewEncoder(w).Encode(platforms)
}

func deletePlatformLocal(platform Platform) {
	var programs []Program
	db.Model(&platform).Association("Programs").Find(&programs)
	for _, s := range programs {
		deleteProgramLocal(s)
	}
	db.Unscoped().Delete(&platform)
}

// Dumps all programs associated with this platform
func getAssociatedPrograms(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var platform Platform
	var programs []Program
	vars := mux.Vars(r)
	db.Where("ID = ?", vars["id"]).Find(&platform)
	db.Model(&platform).Association("Programs").Find(&programs)
	json.NewEncoder(w).Encode(&programs)

}
