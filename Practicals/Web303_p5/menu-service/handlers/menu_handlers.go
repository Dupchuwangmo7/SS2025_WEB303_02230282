package handlers

import (
	"encoding/json"
	"net/http"
	"menu-service/database"
	"menu-service/models"

	"github.com/go-chi/chi/v5"
)

func GetMenu(w http.ResponseWriter, r *http.Request) {
	menuID := chi.URLParam(r, "id")
	var menu models.Menu
	if err := database.DB.Where("id = ?", menuID).First(&menu).Error; err != nil {
		http.Error(w, "Menu item not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(menu)
}

func CreateMenu(w http.ResponseWriter, r *http.Request) {
	var menuData models.Menu
	if err := json.NewDecoder(r.Body).Decode(&menuData); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	result := database.DB.Create(&menuData)
	if result.Error != nil {
		http.Error(w, "Failed to create menu: "+result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(menuData)
}

// func GetMenuItems(w http.ResponseWriter, r *http.Request) {
// 	var items []models.MenuItem
// 	if err := database.DB.Find(&items).Error; err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(items)
// }

// func CreateMenuItem(w http.ResponseWriter, r *http.Request) {
// 	var item models.MenuItem
// 	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	if err := database.DB.Create(&item).Error; err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(item)
// }

// func GetMenuItem(w http.ResponseWriter, r *http.Request) {
// 	id := chi.URLParam(r, "id")

// 	var item models.MenuItem
// 	if err := database.DB.First(&item, id).Error; err != nil {
// 		http.Error(w, "Menu item not found", http.StatusNotFound)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(item)
// }
