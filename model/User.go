package model

import (
	"fmt"
	"os"
)

type User struct {
	UserID      string `gorm:"primary_key; size:50;"`
	SearchIndex string `gorm:"primary_key; size:50;"`
	Handle      bool   `gorm:"not null;"`
}

func CreateUserTable() {
	dbname := fmt.Sprintf("host=%s user=%s dbname=%s  password=%s", os.Getenv("HOST"), os.Getenv("DBUSER"), os.Getenv("DBNAME"), os.Getenv("PASSWORD"))
	ConnectDataBase(dbname)
	if DB.HasTable(&User{}) {
		DB.DropTable("users")
	}
	DB.CreateTable(&User{})
}
