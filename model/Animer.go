package model

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type ACG struct {
	Image    string `gorm:"size:150;"`                       // 首頁影像圖片網址
	TaiName  string `gorm:"primary_key; size:55; not null;"` // 動畫台灣翻譯名稱
	JapName  string `gorm:"size:55; not null;"`              // 動畫日文原名
	Class    string `gorm:"size:5;"`                         // 動畫種類(電影or番)
	Premiere string `gorm:"size30;"`                         // 首播時間
	Author   string `gorm:"size:15;"`                        // 原著作者
	Director string `gorm:"size:15;"`                        // 導演監督
	Firm     string `gorm:"size:10;"`                        // 製作廠商
	Agent    string `gorm:"size:10;"`                        // 台灣代理
	Website  string `gorm:"size:150;"`                       // 官方網站
}

func CreateACGTable() {
	dname := fmt.Sprintf("%s:%s@(%s)/Animers?charset=utf8&parseTime=True&loc=Local", os.Getenv("account"), os.Getenv("password"), os.Getenv("IP"))
	fmt.Println(dname)
	ConnectDataBase(dname)
	DB.CreateTable(&ACG{})
}

// 建立動漫的資料庫
// https://acg.gamer.com.tw/index.php?page=2&p=ANIME&t=1&tnum=5406
func CrewAnimerInfo() {
	db := fmt.Sprintf("%s:%s@(%s)/Animers?charset=utf8&parseTime=True&loc=Local", os.Getenv("account"), os.Getenv("password"), os.Getenv("IP"))
	ConnectDataBase(db)

	// wg1 := new(sync.WaitGroup)
	// wg1.Add(404)

	for i := 1; i <= 404; i++ {
		url := fmt.Sprintf("https://acg.gamer.com.tw/index.php?page=%d&p=ANIME&t=1&tnum=5406", i)
		// go func() {
		// 	urls := CrewSinglePage(url, wg1)
		// 	CrewEachAnime(urls)
		// }()
		urls := CrewSinglePage(url)
		CrewEachAnime(urls)

		// 每爬一頁睡十二秒
		time.Sleep(6 * time.Second)
	}
	// wg1.Wait()

}

func CrewSinglePage(url string) []string {
	var urls []string
	// defer wg.Done()

	dom, _ := getDecument(url)
	dom.Find("div.ACG-mainbox1>div.ACG-mainbox2").Each(func(idx int, s *goquery.Selection) {
		// 把每頁的每部子動畫過濾出來
		ul, _ := s.Find("h1.ACG-maintitle>a").First().Attr("href")
		urls = append(urls, "https:"+ul)
	})

	return urls
}

func CrewEachAnime(urls []string) {
	for _, url := range urls {
		var acg ACG
		dom, _ := getDecument(url)
		s := dom.Find("div.ACG-mster_box1").First()
		acg.Image, _ = s.Find("img").Attr("src")
		acg.TaiName = s.Find("h1").First().Text()
		acg.JapName = s.Find("h2").First().Text()
		acg.Class = s.Find("ul.ACG-box1listA>li:contains(播映方式)").First().Text()
		acg.Premiere = s.Find("ul.ACG-box1listA>li:contains(當地首播)").First().Text()
		s.Find("ul.ACG-box1listB>li").Each(func(idx int, ss *goquery.Selection) {
			switch idx {
			case 0:
				acg.Author = ss.Text()
			case 1:
				acg.Director = ss.Text()
			case 2:
				acg.Firm = ss.Text()
			case 3:
				acg.Agent = ss.Text()
			case 4:
				acg.Website = ss.Find("a").Text()
			}
		})
		log.Println(acg)
		time.Sleep(time.Second * 1)
		DB.Create(&acg)
	}
}

func CrewEachAnimeTest(url string) {
	var acg ACG

	dom, _ := getDecument(url)
	s := dom.Find("div.ACG-mster_box1").First()
	acg.Image, _ = s.Find("img").Attr("src")
	acg.TaiName = s.Find("h1").First().Text()
	acg.JapName = s.Find("h2").First().Text()
	acg.Class = s.Find("ul.ACG-box1listA>li:contains(播映方式)").First().Text()
	acg.Premiere = strings.Split(s.Find("ul.ACG-box1listA>li:contains(當地首播)").First().Text(), "：")[1]
	s.Find("ul.ACG-box1listB>li").Each(func(idx int, ss *goquery.Selection) {
		switch idx {
		case 0:
			acg.Author = ss.Text()
		case 1:
			acg.Director = ss.Text()
		case 2:
			acg.Firm = ss.Text()
		case 3:
			acg.Agent = ss.Text()
		case 4:
			acg.Website = ss.Find("a").Text()
		}
	})

	fmt.Println(acg)
}

func CrewSinglePageTest(url string) {
	dom, _ := getDecument(url)
	dom.Find("div.ACG-mainbox1>div.ACG-mainbox2").Each(func(idx int, s *goquery.Selection) {
		acgurl, _ := s.Find("h1.ACG-maintitle>a").First().Attr("href")
		fmt.Println("https:" + acgurl)
	})
}

// 用Get的方式取得指定網址的html文檔, 並且轉換成goquery用來檢索的strcut
func getDecument(url string) (*goquery.Document, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		log.Printf("錯誤, 請確認您輸入的網址是否正確, 錯誤網址為: %s\n", url)
		return nil, nil
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		fmt.Println("Error, Status code is ", res.StatusCode)
		return nil, errors.New("Status code is not 200!")
	}
	bodyByte, _ := ioutil.ReadAll(res.Body)

	dom, err := goquery.NewDocumentFromReader(bytes.NewReader(bodyByte))
	if err != nil {
		return nil, err
	}
	return dom, nil
}
