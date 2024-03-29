package api

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/davidleitw/AnimeBot/model"
	"github.com/line/line-bot-sdk-go/linebot"
)

func HandleEventTypePostback(event *linebot.Event, bot *linebot.Client) {
	userID := event.Source.UserID
	data := event.Postback.Data
	search, action := handlePostbackData(data)
	switch action {
	case "help":
		_, err := bot.ReplyMessage(
			event.ReplyToken,
			linebot.NewTextMessage(
				FirstHelpMessage,
			),
		).Do()
		if err != nil {
			log.Println("!help message error = ", err)
		}
	case "add":
		// 新增作品至收藏清單
		var users []model.User
		model.DB.Where("user_id = ?", event.Source.UserID).Find(&users)
		if len(users) >= 50 {
			_, err := bot.ReplyMessage(
				event.ReplyToken,
				linebot.NewTextMessage("不好意思，anime bot最多只能在收藏清單放入50部作品喔。\n如果您要再新增資料的話請先把已經觀賞完的作品從清單移除， 謝謝您的配合！"),
			).Do()
			if err != nil {
				log.Println("over added area error = ", err)
			}
		} else {
			err := handleAddItem(userID, search)
			// 如果新增資料沒有錯誤, 回覆新增成功訊息
			if err == nil {
				var anime model.ACG
				model.DB.Where("search_index = ?", search).First(&anime)
				replyMessage := fmt.Sprintf("已成功將%s新增至收藏清單", anime.TaiName)

				_, replyErr := bot.ReplyMessage(
					event.ReplyToken,
					linebot.NewTextMessage(replyMessage),
				).Do()
				// 發送新增成功訊息錯誤時會跳到下面這行
				if replyErr != nil {
					log.Println("Add data result show error = ", replyErr)
				}
			}
		}
	case "delete":
		// 刪除特定項目
		err := handleDeleteItem(userID, search)
		if err == nil {
			var users []model.User
			model.DB.Where("user_id = ?", event.Source.UserID).Find(&users)
			// 刪除完把清單資料再次show出來
			if len(users) != 0 {
				handleUserlist(users, bot, event.ReplyToken)
			}
		}
	case "show":
		// 顯示當前使用者收藏名單
		var users []model.User
		model.DB.Where("user_id = ?", event.Source.UserID).Find(&users)
		handleUserlist(users, bot, event.ReplyToken)

	case "recommand":
		// 新番推薦
		handleRecommand(bot, event.ReplyToken)
	}

	log.Println("user = ", userID, ", search = ", search, ", action = ", action)
}

// Newanimes
func handleRecommand(bot *linebot.Client, token string) {
	var animes model.NewAnimes
	model.DB.Find(&animes)
	sort.Sort(animes)
	animesSubset := animes[:10]
	flex := buildNewAnimesList(animesSubset)
	_, err := bot.ReplyMessage(
		token,
		linebot.NewTextMessage("2020年夏季新番推薦如下:"),
		linebot.NewFlexMessage("新番推薦", flex),
	).Do()
	if err != nil {
		log.Println("New anime error = ", err)
	}
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
	user.Handle = false
	err := model.DB.Create(&user).Error
	return err
}

// 刪除指定項目
func handleDeleteItem(userID, search string) error {
	var user model.User
	user.UserID = userID
	user.SearchIndex = search
	err := model.DB.Delete(&user).Error
	return err
}

func handleUserlist(users []model.User, bot *linebot.Client, token string) {
	l := len(users)
	switch {
	case l == 0:
		_, err := bot.ReplyMessage(
			token,
			linebot.NewTextMessage("目前您的清單還沒有資料"),
		).Do()
		if err != nil {
			log.Println("show function empty error message!")
		}
	case l <= 10:
		flex := buildUserFavoriteList(users)
		_, err := bot.ReplyMessage(
			token,
			linebot.NewFlexMessage("收集清單 page 1", flex),
		).Do()
		if err != nil {
			log.Println("Show list error = ", err)
		}
	case l > 10 && l <= 20:
		flex1 := buildUserFavoriteList(users[:10])
		flex2 := buildUserFavoriteList(users[10:l])
		_, err := bot.ReplyMessage(
			token,
			linebot.NewFlexMessage("收集清單 page 1", flex1),
			linebot.NewFlexMessage("收集清單 page 2", flex2),
		).Do()
		if err != nil {
			log.Println("Show list error = ", err)
		}
	case l > 20 && l <= 30:
		flex1 := buildUserFavoriteList(users[:10])
		flex2 := buildUserFavoriteList(users[10:20])
		flex3 := buildUserFavoriteList(users[20:l])
		_, err := bot.ReplyMessage(
			token,
			linebot.NewFlexMessage("收集清單 page 1", flex1),
			linebot.NewFlexMessage("收集清單 page 2", flex2),
			linebot.NewFlexMessage("收集清單 page 3", flex3),
		).Do()
		if err != nil {
			log.Println("Show list error = ", err)
		}
	case l > 30 && l <= 40:
		flex1 := buildUserFavoriteList(users[:10])
		flex2 := buildUserFavoriteList(users[10:20])
		flex3 := buildUserFavoriteList(users[20:30])
		flex4 := buildUserFavoriteList(users[30:l])
		_, err := bot.ReplyMessage(
			token,
			linebot.NewFlexMessage("收集清單 page 1", flex1),
			linebot.NewFlexMessage("收集清單 page 2", flex2),
			linebot.NewFlexMessage("收集清單 page 3", flex3),
			linebot.NewFlexMessage("收集清單 page 4", flex4),
		).Do()
		if err != nil {
			log.Println("Show list error = ", err)
		}
	case l > 40 && l <= 50:
		flex1 := buildUserFavoriteList(users[:10])
		flex2 := buildUserFavoriteList(users[10:20])
		flex3 := buildUserFavoriteList(users[20:30])
		flex4 := buildUserFavoriteList(users[30:40])
		flex5 := buildUserFavoriteList(users[40:l])
		_, err := bot.ReplyMessage(
			token,
			linebot.NewFlexMessage("收集清單 page 1", flex1),
			linebot.NewFlexMessage("收集清單 page 2", flex2),
			linebot.NewFlexMessage("收集清單 page 3", flex3),
			linebot.NewFlexMessage("收集清單 page 4", flex4),
			linebot.NewFlexMessage("收集清單 page 5", flex5),
		).Do()
		if err != nil {
			log.Println("Show list error = ", err)
		}
	}
}

// Newanimes
func buildNewAnimesList(animes model.NewAnimes) *linebot.CarouselContainer {
	container := &linebot.CarouselContainer{
		Type:     linebot.FlexContainerTypeCarousel,
		Contents: buildBubblesWithNewAnime(animes),
	}
	return container
}

// Newanimes
func buildBubblesWithNewAnime(animes model.NewAnimes) []*linebot.BubbleContainer {
	var containers []*linebot.BubbleContainer
	for _, anime := range animes {
		log.Println("ok")
		containers = append(containers, buildBubbleWithNewAnime(anime))
	}
	return containers
}

// Newanimes
func buildBubbleWithNewAnime(anime model.NewAnime) *linebot.BubbleContainer {
	contain := &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Hero: &linebot.ImageComponent{
			Type:       linebot.FlexComponentTypeImage,
			URL:        anime.ImageSrc,
			Size:       linebot.FlexImageSizeTypeFull,
			AspectMode: linebot.FlexImageAspectModeTypeCover, // 有可能的錯誤1
		},
		Body: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeBaseline,
			Contents: []linebot.FlexComponent{
				&linebot.IconComponent{
					Type: linebot.FlexComponentTypeIcon,
					URL:  "https://img.icons8.com/officel/2x/fire-element.png",
				},
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   anime.TaiName,
					Margin: linebot.FlexComponentMarginTypeMd,
					Size:   linebot.FlexTextSizeTypeMd,
					Weight: linebot.FlexTextWeightTypeBold,
					Color:  "#f7af31",
				},
			},
		},
		Footer: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.ButtonComponent{
					Type:  linebot.FlexComponentTypeButton,
					Style: linebot.FlexButtonStyleTypeLink,
					Color: "#f7af31",
					Action: &linebot.PostbackAction{
						Label: "加入收藏清單",
						Data:  anime.SearchIndex + "&action=add",
					},
					Margin: linebot.FlexComponentMarginTypeXxl,
				},
				&linebot.ButtonComponent{
					Type:  linebot.FlexComponentTypeButton,
					Style: linebot.FlexButtonStyleTypeLink,
					Color: "#f7af31",
					Action: &linebot.URIAction{
						Label: "作品詳細資料",
						URI:   fmt.Sprintf("https://acg.gamer.com.tw/acgDetail.php?s=%s", anime.SearchIndex),
					},
				},
			},
		},
	}
	return contain
}

// build特定使用者清單
func buildUserFavoriteList(users []model.User) *linebot.CarouselContainer {
	container := &linebot.CarouselContainer{
		Type:     linebot.FlexContainerTypeCarousel,
		Contents: buildBubblesWithNewAnimeforlist(users),
	}
	return container
}

// flex message 集合
func buildBubblesWithNewAnimeforlist(users []model.User) []*linebot.BubbleContainer {
	var containers []*linebot.BubbleContainer
	for _, user := range users {
		var anime model.ACG
		search_index := user.SearchIndex
		model.DB.Where("search_index = ?", search_index).First(&anime)
		model.VerifyAnime(&anime)
		containers = append(containers, buildBubbleWithAnimeForList(anime))
	}
	return containers
}

// 單個flex message
func buildBubbleWithAnimeForList(anime model.ACG) *linebot.BubbleContainer {
	container := &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Hero: &linebot.ImageComponent{
			Type:       linebot.FlexComponentTypeImage,
			URL:        anime.Image,
			Size:       linebot.FlexImageSizeTypeFull,
			AspectMode: linebot.FlexImageAspectModeTypeCover, // 有可能的錯誤1
		},
		Body: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeBaseline,
			Contents: []linebot.FlexComponent{
				&linebot.IconComponent{
					Type: linebot.FlexComponentTypeIcon,
					URL:  "https://img.icons8.com/ios-filled/2x/love-book.png",
				},
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   anime.TaiName,
					Margin: linebot.FlexComponentMarginTypeMd,
					Size:   linebot.FlexTextSizeTypeMd,
					Weight: linebot.FlexTextWeightTypeBold,
					Color:  "#f7af31",
				},
			},
		},
		Footer: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.ButtonComponent{
					Type:  linebot.FlexComponentTypeButton,
					Style: linebot.FlexButtonStyleTypeLink,
					Color: "#f7af31",
					Action: &linebot.URIAction{
						Label: "作品詳細資料",
						URI:   fmt.Sprintf("https://acg.gamer.com.tw/acgDetail.php?s=%s", anime.SearchIndex),
					},
					Margin: linebot.FlexComponentMarginTypeXxl,
				},
				&linebot.ButtonComponent{
					Type:  linebot.FlexComponentTypeButton,
					Style: linebot.FlexButtonStyleTypeLink,
					Color: "#f7af31",
					Action: &linebot.PostbackAction{
						Label:       "移除",
						Data:        anime.SearchIndex + "&action=delete",
						DisplayText: "移除此作品",
					},
				},
			},
		},
	}
	return container
}
