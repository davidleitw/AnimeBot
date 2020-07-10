package main

import (
	"github.com/davidleitw/AnimeBot/server"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	linebot := server.AnimeBotServer()
	linebot.Run()
}
