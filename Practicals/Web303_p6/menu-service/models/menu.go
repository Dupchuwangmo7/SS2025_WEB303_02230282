package models

import "gorm.io/gorm"

// Menu represents a café menu collection with related items
type Menu struct {
	gorm.Model
	Name        string     `json:"name"`
	Description string     `json:"description"`
	MenuItems   []MenuItem `json:"menu_items" gorm:"references:MenuID"`
}

// MenuItem represents an individual item available on the café menu
type MenuItem struct {
	gorm.Model
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}
