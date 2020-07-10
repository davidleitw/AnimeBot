package model

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/joho/godotenv/autoload"
)

var DB *gorm.DB

func ConnectDataBase(dbname string) {
	db, _ := gorm.Open("postgres", dbname)

	fmt.Println("Connect database successfully!")
	DB = db
	migration()
}

func migration() {
	DB.AutoMigrate(&ACG{})
}
