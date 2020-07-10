package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/davidleitw/AnimeBot/api"
	"github.com/davidleitw/AnimeBot/model"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
)

var bot *linebot.Client

func AnimeBotServer() *gin.Engine {
	server := gin.Default()
	dbname := fmt.Sprintf("host=%s user=%s dbname=%s  password=%s", os.Getenv("HOST"), os.Getenv("DBUSER"), os.Getenv("DBNAME"), os.Getenv("PASSWORD"))
	model.ConnectDataBase(dbname)

	var err error
	bot, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	if err == nil {
		log.Println("line bot linking successfully!")
	}
	server.POST("/callback", callbackHandler)
	return server
}

func callbackHandler(ctx *gin.Context) {
	events, err := bot.ParseRequest(ctx.Request)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			// 400
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			// 500
			ctx.JSON(http.StatusInternalServerError, err)
		}
	}

	for _, event := range events {
		switch event.Type {
		// 加入好友的時候會觸發的部份
		case linebot.EventTypeFollow:
			api.HandleEventTypeFollow(event, bot)
		// 文字訊息的部份
		case linebot.EventTypeMessage:
			api.HandleEventTypeMessage(event, bot)
		// Postback觸發
		case linebot.EventTypePostback:
			api.HandleEventTypePostback(event, bot)
		// 封鎖linebot的時候
		case linebot.EventTypeUnfollow:
			api.HandleEventTypeUnfollow(event, bot)
		}
	}

}
