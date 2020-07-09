package main

import (
	"github.com/davidleitw/AnimeBot/server"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	linebot := server.AnimeBotServer()
	linebot.Run()
}
