package api

import (
	"log"
	"strings"

	"github.com/davidleitw/AnimeBot/model"

	"github.com/line/line-bot-sdk-go/linebot"
)

func HandleEventTypeMessage(event *linebot.Event, bot *linebot.Client) {
	switch message := event.Message.(type) {
	case *linebot.TextMessage:
		log.Println("Input text: ", message.Text)
		if message.Text == "!help" || message.Text == "-h" || message.Text == "-help" {
			// 功能講解
			log.Println("help area!")
			_, err := bot.ReplyMessage(
				event.ReplyToken,
				linebot.NewTextMessage("Help message: Help!"),
			).Do()
			if err != nil {
				log.Println("!help message error = ", err)
			}
		} else if message.Text == "測試" {
			animes := model.SearchAnimeInfoWithKey("刀劍")
			flex := buildFlexContainerTypeCarousel(animes)
			_, err := bot.ReplyMessage(
				event.ReplyToken,
				linebot.NewFlexMessage("Flex1", flex),
				linebot.NewFlexMessage("Flex2", flex),
			).Do()
			if err != nil {
				log.Println("Testing error = ", err)
			}
		} else if (message.Text[0] == '@' || message.Text[0] == '!') && len(message.Text) >= 2 {
			// 搜尋單一動漫
			animes := model.SearchAnimeInfoWithKey(message.Text[1:])
			if len(animes) > 0 {
				flex := buildFlexContainerTypeCarousel(animes)
				_, err := bot.ReplyMessage(
					event.ReplyToken,
					linebot.NewFlexMessage("flex", flex),
				).Do()
				if err != nil {
					log.Println("Send search response error = ", err)
				}
			} else {
				_, err := bot.ReplyMessage(
					event.ReplyToken,
					linebot.NewTextMessage("對不起, 您輸入的關鍵字無法查詢到結果, 請確認輸入的文字是否正確"),
				).Do()
				if err != nil {
					log.Println("search zero statment error!")
				}
			}

		} else if strings.Contains(message.Text, "https") {
			// 以巴哈姆特網址查詢
			log.Println("https area!")
			anime, err := model.SearchAnimeInfoWithindex(message.Text)
			if err != nil {
				// 沒有搜尋到
				_, err := bot.ReplyMessage(
					event.ReplyToken,
					linebot.NewTextMessage("對不起, 您輸入的網址無法查詢到結果, 請確認輸入的網址是否正確"),
				).Do()
				if err != nil {
					log.Println("search zero statment error!")
				}
			} else {
				flex := buildFlexContainerTypeCarouselSingle(anime)
				_, err := bot.ReplyMessage(
					event.ReplyToken,
					linebot.NewFlexMessage("flex", flex),
				).Do()
				if err != nil {
					log.Println("Send search response error = ", err)
				}
			}

		} else {
			log.Println("else area!")
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
func buildFlexContainerTypeCarouselSingle(anime model.ACG) *linebot.CarouselContainer {
	container := &linebot.CarouselContainer{
		Type: linebot.FlexContainerTypeCarousel,
		Contents: []*linebot.BubbleContainer{
			buildFlexMessageWithAnime(anime),
		},
	}
	return container
}

func buildFlexContainerTypeCarousel(animes []model.ACG) *linebot.CarouselContainer {
	container := &linebot.CarouselContainer{
		Type:     linebot.FlexContainerTypeCarousel,
		Contents: buildFlexContainersTypeBubble(animes),
	}
	return container
}

func buildFlexContainersTypeBubble(animes []model.ACG) []*linebot.BubbleContainer {
	var containers []*linebot.BubbleContainer
	for _, anime := range animes {
		containers = append(containers, buildFlexMessageWithAnime(anime))
	}
	return containers
}

func buildFlexMessageWithAnime(anime model.ACG) *linebot.BubbleContainer {
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
						Label: "添加至欲觀看清單",
						Data:  anime.SearchIndex + "&action=add", // 添加指定的動漫所需要的編號
					},
					Margin: linebot.FlexComponentMarginTypeXxl,
				},
				&linebot.ButtonComponent{
					Type:  linebot.FlexComponentTypeButton,
					Style: linebot.FlexButtonStyleTypePrimary,
					Color: "#f7af31",
					Action: &linebot.MessageAction{
						Label: "按鈕2",
						Text:  "按鈕2測試",
					},
					Margin: linebot.FlexComponentMarginTypeXxl,
				},
			},
		},
	}
	return container
}

func replyFlexMessageTest(animeName string) *linebot.BubbleContainer {
	container := &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Hero: &linebot.ImageComponent{
			URL: "https://p2.bahamut.com.tw/B/ACG/c/96/0000110596.JPG",
			//Size: linebot.FlexImageSizeTypeFull,
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
							Text:   "刀劍神域 Alicization War of Underworld -THE LAST SEASON-",
							Wrap:   true,
							Weight: linebot.FlexTextWeightTypeBold,
							Size:   linebot.FlexTextSizeTypeXl,
							Margin: linebot.FlexComponentMarginTypeMd,
							Color:  "#f7af31",
						},
						&linebot.TextComponent{
							Type: linebot.FlexComponentTypeText,
							Text: "ソードアート・オンライン アリシゼーション War of Underworld -THE LAST SEASON-",
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
									Text:  "2020-04-25",
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
									Text:  "川原礫",
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
									Text:  "A-1 Pictures",
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
									Text:  "https://sao-alicization.net/",
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
					Action: &linebot.MessageAction{
						Label: "按鈕1",
						Text:  "按鈕1測試",
					},
					Margin: linebot.FlexComponentMarginTypeXxl,
				},
				&linebot.ButtonComponent{
					Type:  linebot.FlexComponentTypeButton,
					Style: linebot.FlexButtonStyleTypePrimary,
					Color: "#f7af31",
					Action: &linebot.MessageAction{
						Label: "按鈕2",
						Text:  "按鈕2測試",
					},
					Margin: linebot.FlexComponentMarginTypeXxl,
				},
			},
		},
	}
	return container
}
