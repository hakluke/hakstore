package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// User holds details for user logins, namely the API key
type User struct {
	gorm.Model
	ID  string    `json:"id" gorm:"PrimaryKey"`
	Key uuid.UUID `json:"key" gorm:"type:uuid;primary_key;"`
}

// Define authentication middleware struct
type authenticationMiddleware struct {
	keyUsers map[string]string
}

// Get all of the users out of the database and store them in memory for fast access on every request
func (amw *authenticationMiddleware) Populate() {
	var users []User
	db.Find(&users)
	for _, user := range users {
		amw.keyUsers[user.Key.String()] = user.ID
	}
}

// Middleware function, which will be called for each request
func (amw *authenticationMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("X-API-Key")

		// if the user exists, continue the HTTP serve, otherwise return a 403
		if _, found := amw.keyUsers[key]; found { // If I need to get the actual user at any point, it's the _ that is returned here
			// Pass down the request to the next middleware (or final handler)
			next.ServeHTTP(w, r)
		} else {
			// Write an error and stop the handler chain
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	})
}

// BeforeCreate will set a UUID rather than numeric ID.
func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	uuid := uuid.NewV4()
	user.Key = uuid
	return nil
}

// Get all Users
func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var users []User
	db.Find(&users)
	json.NewEncoder(w).Encode(users)
}

// Get a specific user
func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user User
	vars := mux.Vars(r)
	db.Where("ID = ?", vars["id"]).Find(&user)
	json.NewEncoder(w).Encode(&user)
}

// Creates new users, accepts batches
func createUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var users []User
	_ = json.NewDecoder(r.Body).Decode(&users)
	db.Create(&users)
	json.NewEncoder(w).Encode(users)
}

// Updates a user
func updateUser(w http.ResponseWriter, r *http.Request) {
	// not necessary because createUser() uses UPSERT
}

// Deletes a user
func deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var user User
	db.Where("ID = ?", id).Find(&user)
	db.Unscoped().Delete(&user)
	var users []User
	db.Find(&users)
	json.NewEncoder(w).Encode(users)
}
