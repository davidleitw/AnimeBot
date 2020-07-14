package api

import (
	"log"

	"github.com/line/line-bot-sdk-go/linebot"
)

func HandleEventTypeFollow(event *linebot.Event, bot *linebot.Client) {
	_, err := bot.ReplyMessage(
		event.ReplyToken,
		linebot.NewTextMessage(helpMessage),
	).Do()
	if err != nil {
		log.Println("!help message error = ", err)
	}
}
