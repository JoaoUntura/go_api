package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DbOpen() {

	dsn := "postgresql://postgres:jo@ovitor221203@db.shknwijghijqodznrgdt.supabase.co:5432/postgres"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	DB = db

}
