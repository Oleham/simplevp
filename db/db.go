package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var sVPDB *gorm.DB = OpenDatabase()

func OpenDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("simplevp.db"), &gorm.Config{})
	if err != nil {
		panic("didnt find db")
	}

	db.AutoMigrate(&Setting{})
	db.AutoMigrate(&Job{})
	db.AutoMigrate(&File{})

	return db
}
