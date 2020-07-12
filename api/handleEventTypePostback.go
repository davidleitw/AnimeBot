package api

import (
	"log"

	"github.com/line/line-bot-sdk-go/linebot"
)

func HandleEventTypePostback(event *linebot.Event, bot *linebot.Client) {
	user := event.Source.UserID
	data := event.Postback.Data
	log.Println("user = ", user, ", data = ", data)
}
