package api

import "github.com/line/line-bot-sdk-go/linebot"

func HandleEventTypeFollow(event *linebot.Event, bot *linebot.Client) {
	switch message := event.Message.(type) {
	case *linebot.TextMessage:
		if message.Text == "!help" {
			bot.ReplyMessage(
				event.ReplyToken,
				linebot.NewTextMessage("Help message: Help!"),
			)
		} else {
			bot.ReplyMessage(
				event.ReplyToken,
				linebot.NewTextMessage("Anime Bot: "+message.Text),
			)
		}
	}
}
