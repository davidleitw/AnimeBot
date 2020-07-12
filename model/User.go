package model

import (
	"fmt"
	"os"
)

type User struct {
	ID          string `gorm:"primary_key;"`
	SearchIndex string `gorm:"size:50;"`
}

func CreateUserTable() {
	dbname := fmt.Sprintf("host=%s user=%s dbname=%s  password=%s", os.Getenv("HOST"), os.Getenv("DBUSER"), os.Getenv("DBNAME"), os.Getenv("PASSWORD"))
	ConnectDataBase(dbname)
	if DB.HasTable(&User{}) {
		DB.DropTable("users")
	}
	DB.CreateTable(&User{})
}
