package api

import (
	"fmt"
	"log"
	"strings"

	"github.com/davidleitw/AnimeBot/model"
	"github.com/line/line-bot-sdk-go/linebot"
)

func HandleEventTypePostback(event *linebot.Event, bot *linebot.Client) {
	userID := event.Source.UserID
	data := event.Postback.Data
	search, action := handlePostbackData(data)
	switch action {
	case "add":
		err := handleAddItem(userID, search)
		// 如果新增資料沒有錯誤, 回覆新增成功訊息
		if err == nil {
			_, replyerr := bot.ReplyMessage(
				event.ReplyToken,
				linebot.NewTextMessage("新增成功!"),
			).Do()
			// 發送新增成功訊息錯誤時會跳到下面這行
			if replyerr != nil {
				log.Println("Add data result show error = ", replyerr)
			}
		}
	case "delete":
		err := handleDeleteItem(userID, search)
		if err != nil {
			_, replyerr := bot.ReplyMessage(
				event.ReplyToken,
				linebot.NewTextMessage("刪除成功!"),
			).Do()
			// 發送新增成功訊息錯誤時會跳到下面這行
			if replyerr != nil {
				log.Println("Add data result show error = ", replyerr)
			}
		}
	case "show":
		var users []model.User
		model.DB.Where("user_id = ?", event.Source.UserID).Find(&users)
		if len(users) == 0 {
			_, err := bot.ReplyMessage(
				event.ReplyToken,
				linebot.NewTextMessage("目前您的清單還沒有資料"),
			).Do()
			if err != nil {
				log.Println("show function empty error message!")
			}
			break
		}
		flex := handleShowlist(users)
		_, err := bot.ReplyMessage(
			event.ReplyToken,
			linebot.NewFlexMessage("flex", flex),
		).Do()
		if err != nil {
			log.Println("Show list error = ", err)
		}
	}
	log.Println("user = ", userID, ", search = ", search, ", action = ", action)
}

// search&action=xxx
func handlePostbackData(data string) (search, action string) {
	_data := strings.Split(data, "&")
	search = _data[0]
	action = strings.Split(_data[1], "=")[1]
	return
}

// 用戶ID + search index 新增一個項目至清單, 不會重複
func handleAddItem(userID, search string) error {
	var user model.User
	user.UserID = userID
	user.SearchIndex = search
	err := model.DB.Create(&user).Error
	return err
}

func handleDeleteItem(userID, search string) error {
	var user model.User
	user.UserID = userID
	user.SearchIndex = search
	err := model.DB.Delete(&user).Error
	return err
}

func handleShowlist(users []model.User) *linebot.CarouselContainer {
	container := &linebot.CarouselContainer{
		Type:     linebot.FlexContainerTypeCarousel,
		Contents: buildFlexContainBubbles(users),
	}
	return container
}

func buildFlexContainBubbles(users []model.User) []*linebot.BubbleContainer {
	var containers []*linebot.BubbleContainer
	for _, user := range users {
		var anime model.ACG
		search_index := user.SearchIndex
		model.DB.Where("search_index = ?", search_index).First(&anime)
		containers = append(containers, buildFlexContainCarouselwithItem(anime))
	}
	return containers

}

func buildFlexContainCarouselwithItem(anime model.ACG) *linebot.BubbleContainer {
	container := &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Hero: &linebot.ImageComponent{
			URL:  anime.Image,
			Size: linebot.FlexImageSizeType5xl,
		},
		Body: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.BoxComponent{
					Type:   linebot.FlexComponentTypeBox,
					Layout: linebot.FlexBoxLayoutTypeVertical,
					Contents: []linebot.FlexComponent{
						&linebot.TextComponent{
							Type:   linebot.FlexComponentTypeText,
							Text:   anime.TaiName,
							Wrap:   true,
							Weight: linebot.FlexTextWeightTypeBold,
							Size:   linebot.FlexTextSizeTypeXl,
							Margin: linebot.FlexComponentMarginTypeMd,
							Color:  "#f7af31",
						},
						&linebot.TextComponent{
							Type: linebot.FlexComponentTypeText,
							Text: anime.JapName,
							Size: linebot.FlexTextSizeTypeXs,
							Wrap: true,
						},
						&linebot.SeparatorComponent{},
						&linebot.SpacerComponent{},
					},
				},
			},
			Margin: linebot.FlexComponentMarginTypeXxl,
		},
		Footer: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.ButtonComponent{
					Type:  linebot.FlexComponentTypeButton,
					Style: linebot.FlexButtonStyleTypePrimary,
					Color: "#f7af31",
					Action: &linebot.URIAction{
						Label: "作品詳細資料",
						URI:   fmt.Sprintf("https://acg.gamer.com.tw/acgDetail.php?s=%s", anime.SearchIndex),
					},
					Margin: linebot.FlexComponentMarginTypeXl,
				},
				&linebot.ButtonComponent{
					Type:  linebot.FlexComponentTypeButton,
					Style: linebot.FlexButtonStyleTypePrimary,
					Color: "#f7af31",
					Action: &linebot.PostbackAction{
						Label:       "移除",
						Data:        anime.SearchIndex + "&action=delete",
						DisplayText: "移除此作品",
					},
					Margin: linebot.FlexComponentMarginTypeXl,
				},
			},
		},
	}

	return container
}
