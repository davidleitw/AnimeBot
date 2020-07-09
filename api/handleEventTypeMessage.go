package api

import (
	"log"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
)

func HandleEventTypeMessage(event *linebot.Event, bot *linebot.Client) {
	switch message := event.Message.(type) {
	case *linebot.TextMessage:
		if message.Text == "!help" || message.Text == "-h" || message.Text == "-help" {
			// 功能講解
			_, err := bot.ReplyMessage(
				event.ReplyToken,
				linebot.NewTextMessage("Help message: Help!"),
			).Do()
			if err != nil {
				log.Println("!help message error = ", err)
			}
		} else if message.Text[0] == '@' {
			// 搜尋單一動漫
			name := strings.Split(message.Text, "@")[1]
			_, err := bot.ReplyMessage(
				event.ReplyToken,
				linebot.NewTextMessage("作品名稱: "+name),
			).Do()
			if err != nil {
				log.Println("Send search response error = ", err)
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
