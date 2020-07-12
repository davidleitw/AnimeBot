package api

import (
	"log"
	"strings"

	"github.com/davidleitw/AnimeBot/model"
	"github.com/line/line-bot-sdk-go/linebot"
)

func HandleEventTypePostback(event *linebot.Event, bot *linebot.Client) {
	user := event.Source.UserID
	data := event.Postback.Data
	search, action := handlePostbackData(data)
	switch action {
	case "add":
		handleAddItem(user, search)
	case "delete":
		handleDeleteItem(user, search)
	case "show":
		handleShowList(user, bot)
	}
	log.Println("user = ", user, ", search = ", search, ", action = ", action)
}

// search&action=xxx
func handlePostbackData(data string) (search, action string) {
	_data := strings.Split(data, "&")
	search = _data[0]
	action = strings.Split(_data[1], "=")[1]
	return
}

func handleAddItem(userID, search string) {
	var user model.User
	user.UserID = userID
	user.SearchIndex = search
	model.DB.Create(&user)
}

func handleDeleteItem(userID, search string) {

}

func handleShowList(userID string, bot *linebot.Client) {

}
