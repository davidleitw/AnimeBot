package main

import (
	"github.com/davidleitw/AnimeBot/model"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	linebot := server.AnimeBotServer()
	err := linebot.Run()
	if err != nil {
		fmt.Println("Server running error: ", err)
	}

	// model.SearchWithTimeTest()
}
