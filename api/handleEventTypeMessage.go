package api

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/davidleitw/AnimeBot/model"

	"github.com/line/line-bot-sdk-go/linebot"
)

// 針對初次使用或者輸入help的人發送的message
const FirstHelpMessage = "歡迎使用Anime Bot服務!\n" + "此服務可以提供動漫作品查詢，並將其列入收藏清單內\n\n" + "以作品名稱作為關鍵字查詢請輸入@作品名稱\n\n" + "以作者名稱作為關鍵字查詢請輸入!作者名稱\n\n" + "如果上述兩種方法都無法找到指定的作品，可以輸入巴哈姆特動畫資料庫該作品的網址搜尋，即可獲得指定的作品\n" + "ex: https://acg.gamer.com.tw/acgDetail.php?s=110596\n\n" + "查看現在喜好清單內的作品請點擊下方清單按鈕或者輸入@清單\n\n" + "新番推薦請點擊下方的推薦按鈕\n\n" + "參考使用說明請輸入help(大小寫皆可)\n"

// 針對輸入的文字並沒有在功能選單上的使用者所回覆的message
const SencondHelpMessage = "您所輸入的指令並不在預設指令中。\n\n" + "以作品名稱作為關鍵字查詢請輸入@作品名稱\n\n" + "以作者名稱作為關鍵字查詢請輸入!作者名稱\n\n" + "如果上述兩種方法都無法找到指定的作品，可以輸入巴哈姆特動畫資料庫該作品的網址搜尋，即可獲得指定的作品\n" + "ex: https://acg.gamer.com.tw/acgDetail.php?s=110596\n\n" + "查看現在喜好清單內的作品請點擊下方清單按鈕或者輸入@清單\n\n" + "新番推薦請點擊下方的推薦按鈕\n\n" + "參考使用說明請輸入help(大小寫皆可)\n"

func HandleEventTypeMessage(event *linebot.Event, bot *linebot.Client) {
	switch message := event.Message.(type) {
	case *linebot.TextMessage:
		switch Message := message.Text; {
		// 每日一抽, 從人氣高的作品隨機抽出一件
		case Message == "抽":
			anime := handleRandAnime()
			_, err := bot.ReplyMessage(
				event.ReplyToken,
				linebot.NewFlexMessage("抽", buildFlexMessageWithAnime(anime)),
			).Do()
			if err != nil {
				log.Println("send rand anime information area error = ", err)
			}
		// 使用說明
		case strings.EqualFold(Message, "help") || Message == "!help" || Message == "-h" || Message == "-help":
			_, err := bot.ReplyMessage(
				event.ReplyToken,
				linebot.NewTextMessage(FirstHelpMessage),
			).Do()
			if err != nil {
				log.Println("!help message error = ", err)
			}

		// 收藏清單
		case strings.EqualFold(Message, "list") || Message == "@清單" || Message == "清單":
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

		// 以作品名稱查詢
		case Message[0] == '@' && len([]rune(Message)) >= 2:
			animes := model.SearchAnimeInfoWithKey(Message[1:])
			if len(animes) > 0 {
				// 至少有查詢到一個結果
				flex := buildFlexContainerTypeCarousel(animes)
				_, err := bot.ReplyMessage(
					event.ReplyToken,
					linebot.NewFlexMessage("flex", flex),
				).Do()
				if err != nil {
					log.Println("Send search response error = ", err)
				}
			} else {
				// 沒有搜尋到結果
				_, err := bot.ReplyMessage(
					event.ReplyToken,
					linebot.NewTextMessage("對不起, 您輸入的關鍵字無法查詢到結果, 請確認輸入的文字是否正確"),
				).Do()
				if err != nil {
					log.Println("search zero statment error!")
				}
			}

		// 以作者名稱查詢
		case Message[0] == '!' && len([]rune(Message)) >= 2:
			animes := model.SearchAnimeInfoWithAuthor(Message[1:])
			if len(animes) > 0 {
				// 至少有查詢到一個結果
				flex := buildFlexContainerTypeCarousel(animes)
				_, err := bot.ReplyMessage(
					event.ReplyToken,
					linebot.NewFlexMessage("flex", flex),
				).Do()
				if err != nil {
					log.Println("Send search response error = ", err)
				}
			} else {
				// 沒有搜尋到結果
				_, err := bot.ReplyMessage(
					event.ReplyToken,
					linebot.NewTextMessage("對不起, 您輸入的關鍵字無法查詢到結果, 請確認輸入的文字是否正確"),
				).Do()
				if err != nil {
					log.Println("search zero statment error!")
				}
			}

		// 以巴哈姆特網址查詢
		case strings.Contains(Message, "https"):
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

		// 非指令區
		default:
			_, err := bot.ReplyMessage(
				event.ReplyToken,
				linebot.NewTextMessage(SencondHelpMessage),
			).Do()
			if err != nil {
				log.Println("!help message error = ", err)
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

func handleRandAnime() model.ACG {
	var anime model.ACG
	var year string
	rand.Seed(time.Now().Unix())

	// 從 2000年到2020年先選擇一年 0 ~ 20
	ly := rand.Intn(21)
	if ly >= 0 && ly < 10 {
		// 2001
		year = "200" + strconv.Itoa(ly)
	} else {
		// 2010
		year = "20" + strconv.Itoa(ly)
	}

	// 隨機存取符合條件的一部作品
	model.DB.Where("premiere LIKE ?", year+"%").Order("RAND()").Find(&anime)
	return anime
}

func handleRandAnimeTest() {
	dbname := fmt.Sprintf("host=%s user=%s dbname=%s  password=%s", os.Getenv("HOST"), os.Getenv("DBUSER"), os.Getenv("DBNAME"), os.Getenv("PASSWORD"))
	model.ConnectDataBase(dbname)

	for i := 0; i < 10; i++ {
		var anime model.ACG
		var year string
		rand.Seed(time.Now().Unix())
		// 從 2000年到2020年先選擇一年 0 ~ 20

		// 隨機存取符合條件的一部作品
		model.DB.Where("premiere LIKE ?", year+"%").Order("RANDOM()").Find(&anime)

		fmt.Printf("name of anime: %s, time of anime: %s, rand choose year: %s\n", anime.TaiName, anime.Premiere, year)
	}

}
