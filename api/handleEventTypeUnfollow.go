package api

import (
	"github.com/line/line-bot-sdk-go/linebot"
)

func HandleEventTypeUnfollow(event *linebot.Event, bot *linebot.Client) {
	// 封鎖line bot的時候刪除該用戶資料庫的資料
	// userID := event.Source.UserID
	// model.DB.Where("user_id = ?", userID).Delete(model.User{})
	// 為了收集資料 所以暫時不要開啟這個功能
}
