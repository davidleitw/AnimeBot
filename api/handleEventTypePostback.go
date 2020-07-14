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
				helpMessage,
			),
		).Do()
		if err != nil {
			log.Println("!help message error = ", err)
		}
	case "add":
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
				_, replyerr := bot.ReplyMessage(
					event.ReplyToken,
					linebot.NewTextMessage("新增成功!"),
				).Do()
				// 發送新增成功訊息錯誤時會跳到下面這行
				if replyerr != nil {
					log.Println("Add data result show error = ", replyerr)
				}
			}
		}

	case "delete":
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
		var animes model.NewAnimes
		model.DB.Find(&animes)
		sort.Sort(animes)
		animesSubset := animes[:10]
		flex := buildNewAnimeslist(animesSubset)
		log.Println("flex = ", flex)
		//flex := buildFlexContainBubblesWithNewAnimes(animesSubset[0])
		_, err := bot.ReplyMessage(
			event.ReplyToken,
			linebot.NewTextMessage("測試reply是否正常"),
			linebot.NewFlexMessage("新番推薦", flex),
		).Do()
		if err != nil {
			log.Println("New anime error = ", err)
		}
	}

	log.Println("user = ", userID, ", search = ", search, ", action = ", action)
}

// Newanimes
func handleRecommand(bot *linebot.Client, token string) {

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
		flex := buildShowlist(users)
		_, err := bot.ReplyMessage(
			token,
			linebot.NewFlexMessage("收集清單 page 1", flex),
		).Do()
		if err != nil {
			log.Println("Show list error = ", err)
		}
	case l > 10 && l <= 20:
		flex1 := buildShowlist(users[:10])
		flex2 := buildShowlist(users[10:l])
		_, err := bot.ReplyMessage(
			token,
			linebot.NewFlexMessage("收集清單 page 1", flex1),
			linebot.NewFlexMessage("收集清單 page 2", flex2),
		).Do()
		if err != nil {
			log.Println("Show list error = ", err)
		}
	case l > 20 && l <= 30:
		flex1 := buildShowlist(users[:10])
		flex2 := buildShowlist(users[10:20])
		flex3 := buildShowlist(users[20:l])
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
		flex1 := buildShowlist(users[:10])
		flex2 := buildShowlist(users[10:20])
		flex3 := buildShowlist(users[20:30])
		flex4 := buildShowlist(users[30:l])
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
		flex1 := buildShowlist(users[:10])
		flex2 := buildShowlist(users[10:20])
		flex3 := buildShowlist(users[20:30])
		flex4 := buildShowlist(users[30:40])
		flex5 := buildShowlist(users[40:l])
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
func buildNewAnimeslist(animes model.NewAnimes) *linebot.CarouselContainer {
	container := &linebot.CarouselContainer{
		Type:     linebot.FlexContainerTypeCarousel,
		Contents: buildFlexContainBubblesNewAnimes(animes),
	}
	return container
}

// Newanimes
func buildFlexContainBubblesNewAnimes(animes model.NewAnimes) []*linebot.BubbleContainer {
	var containers []*linebot.BubbleContainer
	for _, anime := range animes {
		log.Println("ok")
		containers = append(containers, buildFlexContainBubblesWithNewAnimes(anime))
	}
	return containers
}

// Newanimes
func buildFlexContainBubblesWithNewAnimes(anime model.NewAnime) *linebot.BubbleContainer {
	contain := &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Hero: &linebot.ImageComponent{
			URL:  anime.ImageSrc,
			Size: linebot.FlexImageSizeType5xl,
		},
		Body: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.SeparatorComponent{
					Margin: linebot.FlexComponentMarginTypeXxl,
				},
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
				&linebot.BoxComponent{
					Type:   linebot.FlexComponentTypeBox,
					Layout: linebot.FlexBoxLayoutTypeVertical,
					Contents: []linebot.FlexComponent{
						&linebot.BoxComponent{
							Type:   linebot.FlexComponentTypeBox,
							Layout: linebot.FlexBoxLayoutTypeVertical,
							Contents: []linebot.FlexComponent{
								&linebot.TextComponent{
									Type:  linebot.FlexComponentTypeText,
									Text:  "首播時間",
									Size:  linebot.FlexTextSizeTypeSm,
									Color: "#f7af31",
								},
								&linebot.TextComponent{
									Type:  linebot.FlexComponentTypeText,
									Text:  anime.Premiere,
									Size:  linebot.FlexTextSizeTypeSm,
									Color: "#111111",
									Align: linebot.FlexComponentAlignTypeEnd,
								},
							},
						},
						&linebot.BoxComponent{
							Type:   linebot.FlexComponentTypeBox,
							Layout: linebot.FlexBoxLayoutTypeVertical,
							Contents: []linebot.FlexComponent{
								&linebot.TextComponent{
									Type:  linebot.FlexComponentTypeText,
									Text:  "原著作者",
									Size:  linebot.FlexTextSizeTypeSm,
									Color: "#f7af31",
								},
								&linebot.TextComponent{
									Type:  linebot.FlexComponentTypeText,
									Text:  anime.Author,
									Size:  linebot.FlexTextSizeTypeSm,
									Color: "#111111",
									Align: linebot.FlexComponentAlignTypeEnd,
								},
							},
						},
						&linebot.BoxComponent{
							Type:   linebot.FlexComponentTypeBox,
							Layout: linebot.FlexBoxLayoutTypeVertical,
							Contents: []linebot.FlexComponent{
								&linebot.TextComponent{
									Type:  linebot.FlexComponentTypeText,
									Text:  "作畫公司",
									Size:  linebot.FlexTextSizeTypeSm,
									Color: "#f7af31",
								},
								&linebot.TextComponent{
									Type:  linebot.FlexComponentTypeText,
									Text:  anime.Firm,
									Size:  linebot.FlexTextSizeTypeSm,
									Color: "#111111",
									Align: linebot.FlexComponentAlignTypeEnd,
								},
							},
						},
						&linebot.BoxComponent{
							Type:   linebot.FlexComponentTypeBox,
							Layout: linebot.FlexBoxLayoutTypeVertical,
							Contents: []linebot.FlexComponent{
								&linebot.TextComponent{
									Type:  linebot.FlexComponentTypeText,
									Text:  "官方網站",
									Size:  linebot.FlexTextSizeTypeSm,
									Color: "#f7af31",
								},
								&linebot.TextComponent{
									Type:  linebot.FlexComponentTypeText,
									Text:  anime.Website,
									Size:  linebot.FlexTextSizeTypeSm,
									Color: "#111111",
									Align: linebot.FlexComponentAlignTypeEnd,
								},
							},
						},
						&linebot.SpacerComponent{},
					},
					Margin: linebot.FlexComponentMarginTypeXl,
				},
				&linebot.ButtonComponent{
					Type:  linebot.FlexComponentTypeButton,
					Style: linebot.FlexButtonStyleTypePrimary,
					Color: "#f7af31",
					Action: &linebot.PostbackAction{
						Label:       "添加至欲觀看清單",
						Data:        anime.SearchIndex + "&action=add", // 添加指定的動漫所需要的編號
						DisplayText: "加入清單",
					},
					Margin: linebot.FlexComponentMarginTypeXxl,
				},
				&linebot.ButtonComponent{
					Type:  linebot.FlexComponentTypeButton,
					Style: linebot.FlexButtonStyleTypePrimary,
					Color: "#f7af31",
					Action: &linebot.URIAction{
						Label: "作品詳細資料",
						URI:   fmt.Sprintf("https://acg.gamer.com.tw/acgDetail.php?s=%s", anime.SearchIndex),
					},
					Margin: linebot.FlexComponentMarginTypeXxl,
				},
			},
		},
	}
	return contain
}

// build特定使用者清單
func buildShowlist(users []model.User) *linebot.CarouselContainer {
	container := &linebot.CarouselContainer{
		Type:     linebot.FlexContainerTypeCarousel,
		Contents: buildFlexContainBubbles(users),
	}
	return container
}

// flex message 集合
func buildFlexContainBubbles(users []model.User) []*linebot.BubbleContainer {
	var containers []*linebot.BubbleContainer
	for _, user := range users {
		var anime model.ACG
		search_index := user.SearchIndex
		model.DB.Where("search_index = ?", search_index).First(&anime)
		model.VerifyAnime(&anime)
		containers = append(containers, buildFlexContainCarouselwithItem(anime))
	}
	return containers
}

// 單個flex message
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
