package models

import "gorm.io/gorm"

type Item struct {
	gorm.Model
	Name     string
	Url      string
	Username string
	Password string
	VaultID  uint
}
