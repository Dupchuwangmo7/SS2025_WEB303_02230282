package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type FoodItem struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

var foodItems = []FoodItem{
	{ID: "1", Name: "Espresso", Price: 2.75},
	{ID: "2", Name: "Turkey Sandwich", Price: 5.50},
	{ID: "3", Name: "Blueberry Muffin", Price: 3.50},
	{ID: "4", Name: "Iced Tea", Price: 2.25},
	{ID: "5", Name: "Caesar Salad", Price: 6.00},
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.Get("/items", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(foodItems)
	})

	log.Println("Food Catalog Service starting on port 8080...")
	http.ListenAndServe(":8080", r)
}