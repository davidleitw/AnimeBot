package model

import (
	"fmt"
	"os"

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

func Test() {
	dbname := fmt.Sprintf("host=%s user=%s dbname=%s  password=%s", os.Getenv("HOST"), os.Getenv("DBUSER"), os.Getenv("DBNAME"), os.Getenv("PASSWORD"))
	ConnectDataBase(dbname)

	var animes []ACG
	DB.Find(&animes)
	fmt.Printf("len of animes: %d\n", len(animes))
}
