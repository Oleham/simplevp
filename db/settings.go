// settings.go contain functions for updating and retrieving user settings.
package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ShowSettings() *Setting {

	db, err := gorm.Open(sqlite.Open("simplevp.db"), &gorm.Config{})
	if err != nil {
		panic("didnt find db")
	}

	db.AutoMigrate(&Setting{})

	var current Setting

	db.FirstOrCreate(&current, 1)

	return &current
}

func UpdateSettings(entry *Setting) {
	// Updates the settings
	db, err := gorm.Open(sqlite.Open("simplevp.db"), &gorm.Config{})
	if err != nil {
		panic("didnt find db")
	}

	db.Save(&entry)

}
