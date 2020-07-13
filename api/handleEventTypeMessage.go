package api

import (
	"fmt"
	"log"
	"strings"

	"github.com/davidleitw/AnimeBot/model"

	"github.com/line/line-bot-sdk-go/linebot"
)

func HandleEventTypeMessage(event *linebot.Event, bot *linebot.Client) {
	switch message := event.Message.(type) {
	case *linebot.TextMessage:
		log.Println("Input text: ", message.Text)
		if message.Text == "!help" || message.Text == "-h" || message.Text == "-help" || strings.EqualFold(message.Text, "help") {
			// 功能講解
			log.Println("help area!")
			_, err := bot.ReplyMessage(
				event.ReplyToken,
				linebot.NewTextMessage(`歡迎使用Anime Bot服務!\n
										此服務可以提供動漫作品查詢, 並將其列入喜好清單內\n
										如果想要查詢作品請輸入@作品名稱, 即可跳出搜尋結果\n
										如果以關鍵字搜尋不到可能是因為作品翻譯問題, 可以輸入巴哈姆特動畫資料庫該作品的網址做查詢\n
										ex: https://acg.gamer.com.tw/acgDetail.php?s=110596\n
										如果想查看現在喜好清單內的作品可以點擊下方清單按鈕或者輸入@清單\n`),
			).Do()
			if err != nil {
				log.Println("!help message error = ", err)
			}
		} else if message.Text == "@清單" {
			_, err := bot.ReplyMessage(
				event.ReplyToken,
				linebot.NewFlexMessage("flex", &linebot.BubbleContainer{
					Type: linebot.FlexContainerTypeBubble,
					Footer: &linebot.BoxComponent{
						Type:   linebot.FlexComponentTypeBox,
						Layout: linebot.FlexBoxLayoutTypeVertical,
						Contents: []linebot.FlexComponent{
							&linebot.ButtonComponent{
								Type:  linebot.FlexComponentTypeButton,
								Style: linebot.FlexButtonStyleTypePrimary,
								Color: "#f7af31",
								Action: &linebot.PostbackAction{
									Label: "顯示清單",
									Data:  "000000&action=show",
								},
							},
						},
					},
				}),
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
			if err != nil || anime.IsEmpty() {
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
	return container
}
