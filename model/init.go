package model

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/joho/godotenv/autoload"
)

var DB *gorm.DB

func ConnectDataBase(dbname string) {
	db, _ := gorm.Open("mysql", dbname)

	fmt.Println("Connect database successfully!")
	DB = db
	migration()
}

func migration() {
	DB.Set("gorm:table_options", "charset=utf8mb4").AutoMigrate(&ACG{})
}
