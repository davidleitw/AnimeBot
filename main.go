package main

import (
	"fmt"

	"github.com/davidleitw/AnimeBot/server"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// 待新增 把新番列表的部份高清楚圖移動到資料庫中
	linebot := server.AnimeBotServer()
	err := linebot.Run()
	if err != nil {
		fmt.Println("Server running error: ", err)
	}

}
