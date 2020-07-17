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

func SearchUserInfo(userID string) {
	dbname := fmt.Sprintf("host=%s user=%s dbname=%s  password=%s", os.Getenv("HOST"), os.Getenv("DBUSER"), os.Getenv("DBNAME"), os.Getenv("PASSWORD"))
	ConnectDataBase(dbname)

	var users []User
	DB.Where("user_id = ?", userID).Find(&users)
	for _, user := range users {
		var anime ACG
		search_index := user.SearchIndex
		DB.Where("search_index = ?", search_index).First(&anime)
		fmt.Println(anime.TaiName)
	}
}
