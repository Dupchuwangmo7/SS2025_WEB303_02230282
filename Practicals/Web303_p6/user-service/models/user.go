package models

import "gorm.io/gorm"

// User represents a system user with optional café owner privileges
type User struct {
	gorm.Model
	Name       string `json:"name"`
	Email      string `json:"email" gorm:"unique"`
	IsCafeOwner bool   `json:"is_cafe_owner"`
}
