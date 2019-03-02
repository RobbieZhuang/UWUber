package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Profile struct {
	gorm.Model
	Code string
	Price uint
}

func main() {
	db, err := gorm.Open("sqlite3", "sqlite/user_profile.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&Profile{})

	// Create
	db.Create(&Profile{Code: "L1212", Price: 1000})

	// Read
	var profile Profile
	db.First(&profile, 1) // find product with id 1
	db.First(&profile, "code = ?", "L1212") // find product with code l1212

	// Update - update product's price to 2000
	db.Model(&profile).Update("Price", 2000)
}
