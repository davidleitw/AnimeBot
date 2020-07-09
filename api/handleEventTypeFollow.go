package api

import (
	"log"

	"github.com/line/line-bot-sdk-go/linebot"
)

func HandleEventTypeFollow(event *linebot.Event, bot *linebot.Client) {
	switch message := event.Message.(type) {
	case *linebot.TextMessage:
		if message.Text == "!help" {
			_, err := bot.ReplyMessage(
				event.ReplyToken,
				linebot.NewTextMessage("Help message: Help!"),
			).Do()
			if err != nil {
				log.Println("!help message error = ", err)
			}
		} else {
			_, err := bot.ReplyMessage(
				event.ReplyToken,
				linebot.NewTextMessage("Anime Bot: "+message.Text),
			).Do()
			if err != nil {
				log.Println("Normal text message error = ", err)
			}
		}
	}
}
