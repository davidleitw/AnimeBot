package main

import (
	"fmt"

	"github.com/davidleitw/AnimeBot/server"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	linebot := server.AnimeBotServer()
	err := linebot.Run()
	if err != nil {
		fmt.Println("Server running error: ", err)
	}
	// dbname := fmt.Sprintf("host=%s user=%s dbname=%s  password=%s", os.Getenv("HOST"), os.Getenv("DBUSER"), os.Getenv("DBNAME"), os.Getenv("PASSWORD"))
	// model.ConnectDataBase(dbname)
	// a := model.TestSql("ter")
	// fmt.Println("\n")
	// for _, val := range a {
	// 	fmt.Println(val.JapName, "=> len: ", len(val.JapName))
	// }
	// model.CreateUserTable()
}
