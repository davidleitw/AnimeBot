package api

import (
	"fmt"
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
