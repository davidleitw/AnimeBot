package model

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/joho/godotenv/autoload"
)

var DB *gorm.DB

func ConnectDataBase(dbname string) {
	db, err := gorm.Open("postgres", dbname)
	if err != nil {
		fmt.Println("Connect database error = ", err)
	} else {
		fmt.Println("Connect database successfully!")

	}
	DB = db
	migration()
}

func migration() {
	DB.AutoMigrate(&ACG{})
}
