package api

import (
	"log"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
)

func HandleEventTypeMessage(event *linebot.Event, bot *linebot.Client) {
	switch message := event.Message.(type) {
	case *linebot.TextMessage:
		if message.Text == "!help" || message.Text == "-h" || message.Text == "-help" {
			// 功能講解
			_, err := bot.ReplyMessage(
				event.ReplyToken,
				linebot.NewTextMessage("Help message: Help!"),
			).Do()
			if err != nil {
				log.Println("!help message error = ", err)
			}
		} else if message.Text == "測試" {
			_, err := bot.ReplyMessage(
				event.ReplyToken,
				linebot.NewFlexMessage("Flex", replyFlexMessage("測試")),
			).Do()
			if err != nil {
				log.Println("Testing error = ", err)
			}
		} else if (message.Text[0] == '@' || message.Text[0] == '!') && len(message.Text) >= 2 {
			// 搜尋單一動漫
			split := string(message.Text[0])
			name := strings.Split(message.Text, split)[1]
			_, err := bot.ReplyMessage(
				event.ReplyToken,
				linebot.NewTextMessage("作品名稱: "+name),
			).Do()
			if err != nil {
				log.Println("Send search response error = ", err)
			}
		} else {
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

func replyFlexMessage(animeName string) *linebot.BubbleContainer {
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
					Margin: linebot.FlexComponentMarginTypeNone,
				},
				&linebot.ButtonComponent{
					Type:  linebot.FlexComponentTypeButton,
					Style: linebot.FlexButtonStyleTypePrimary,
					Color: "#f7af31",
					Action: &linebot.MessageAction{
						Label: "按鈕2",
						Text:  "按鈕2測試",
					},
				},
				&linebot.SpacerComponent{},
				&linebot.ButtonComponent{
					Type:  linebot.FlexComponentTypeButton,
					Style: linebot.FlexButtonStyleTypePrimary,
					Color: "#f7af31",
					Action: &linebot.MessageAction{
						Label: "按鈕2",
						Text:  "按鈕2測試",
					},
				},
			},
		},
	}

	return container
}
