package handlers

import (
	"encoding/json"
	"net/http"
	"user-service/database"
	"user-service/models"

	"github.com/go-chi/chi/v5"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var userData models.User
	if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
		http.Error(w, "Invalid user data: "+err.Error(), http.StatusBadRequest)
		return
	}

	result := database.DB.Create(&userData)
	if result.Error != nil {
		http.Error(w, "Failed to create user: "+result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userData)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")

	var user models.User
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		http.Error(w, "User not found with ID: "+userID, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
func GetUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	result := database.DB.Find(&users)
	if result.Error != nil {
		http.Error(w, "Failed to retrieve users: "+result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}json.NewEncoder(w).Encode(users)
}
