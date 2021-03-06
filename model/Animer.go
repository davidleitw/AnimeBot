package model

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type ACG struct {
	SearchIndex string `gorm:"size:50;"`                        // 動漫編號
	Image       string `gorm:"size:150;"`                       // 首頁影像圖片網址
	TaiName     string `gorm:"primary_key; size:90; not null;"` // 動畫台灣翻譯名稱
	JapName     string `gorm:"size:90; not null;"`              // 動畫日文原名
	Class       string `gorm:"size:60;"`                        // 動畫種類(電影or番)
	Premiere    string `gorm:"size:60;"`                        // 首播時間
	Author      string `gorm:"size:60;"`                        // 原著作者
	Director    string `gorm:"size:60;"`                        // 導演監督
	Firm        string `gorm:"size:60;"`                        // 製作廠商
	Agent       string `gorm:"size:60;"`                        // 台灣代理
	Website     string `gorm:"size:150;"`                       // 官方網站
	Popularity  int
}

func (anime ACG) IsEmpty() bool {
	return reflect.DeepEqual(anime, ACG{})
}

// 更新人氣以及還未放入資料庫的作品
func UpdateAnimesInfo() {
	maxPage := 404
	for page := 2; page <= maxPage; page++ {
		pageUrl := fmt.Sprintf("https://acg.gamer.com.tw/index.php?page=%d&p=ANIME&t=1&tnum=5406", page)
		dom, _ := getDecument(pageUrl)

		dom.Find("div.ACG-mainbox1").Each(func(idx int, selection *goquery.Selection) {
			var anime ACG
			// 獲得每個分頁底下動漫的網址
			animeUrl, _ := selection.Find("div.ACG-mainbox2>h1.ACG-maintitle>a").First().Attr("href")
			animeUrl = "https:" + animeUrl
			animePopularity, _ := strconv.Atoi(selection.Find("div.ACG-mainbox4>p.ACG-mainplay>span").First().Text())

			parse, _ := url.Parse(animeUrl)
			query, _ := url.ParseQuery(parse.RawQuery)
			sIndex := query.Get("s")

			err := DB.Where("search_index = ?", sIndex).First(&anime).Error
			// 代表該部作品沒有收錄在資料庫內部
			if err != nil {
				dom, _ := getDecument(animeUrl)
				box := dom.Find("div.ACG-mster_box1").First()
				anime.SearchIndex = sIndex
				anime.Image, _ = box.Find("img").Attr("src")
				anime.TaiName = box.Find("h1").First().Text()
				anime.JapName = box.Find("h2").First().Text()
				anime.Class = CheckColonExist(box.Find("ul.ACG-box1listA>li:contains(播映方式)").First().Text())
				anime.Premiere = CheckColonExist(box.Find("ul.ACG-box1listA>li:contains(當地首播)").First().Text())

				box.Find("ul.ACG-box1listB>li").Each(func(idx int, ss *goquery.Selection) {
					switch idx {
					case 0:
						// 爬取作者欄位
						anime.Author = CheckColonExist(ss.Text())
					case 1:
						// 爬取監督欄位
						anime.Director = CheckColonExist(ss.Text())
					case 2:
						// 爬取製作廠商
						anime.Firm = CheckColonExist(ss.Text())
					case 3:
						// 爬取台灣代理
						anime.Agent = CheckColonExist(ss.Text())
					case 4:
						// 爬取官方網站
						anime.Website = ss.Find("a").Text()
					}
				})

				anime.Popularity = animePopularity
				fmt.Println(anime.TaiName)
				DB.Create(&anime)
			} else {
				// 如果作品存在 更新該作品的人氣
				fmt.Println("update: ", anime.TaiName)
				DB.Model(&anime).Where("search_index = ?", sIndex).Update("popularity", animePopularity)
			}
		})
	}
}

func CreateACGTable() {
	dbname := fmt.Sprintf("host=%s user=%s dbname=%s  password=%s", os.Getenv("HOST"), os.Getenv("DBUSER"), os.Getenv("DBNAME"), os.Getenv("PASSWORD"))
	ConnectDataBase(dbname)
	// 如果要洗掉資料庫重來 把這邊註解刪掉
	// if DB.HasTable(&ACG{}) {
	// 	DB.DropTable("acgs")
	// }
	DB.CreateTable(&ACG{})
}

// 建立動漫的資料庫
// https://acg.gamer.com.tw/index.php?page=2&p=ANIME&t=1&tnum=5406
func CrewAnimerInfo() {
	db := fmt.Sprintf("host=%s user=%s dbname=%s  password=%s", os.Getenv("HOST"), os.Getenv("DBUSER"), os.Getenv("DBNAME"), os.Getenv("PASSWORD"))
	ConnectDataBase(db)

	for i := 1; i <= 404; i++ {
		url := fmt.Sprintf("https://acg.gamer.com.tw/index.php?page=%d&p=ANIME&t=1&tnum=5406", i)
		urls := CrewSinglePage(url)
		CrewEachAnime(urls)
		// 每爬一頁睡十二秒
		time.Sleep(3 * time.Second)
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

func CheckColonExist(str string) string {
	if len(str) == 0 {
		return "nil"
	}
	if val := strings.Split(str, "："); len(val) >= 2 {
		return val[1]
	}
	return str
}

func CrewEachAnime(urls []string) {
	for _, _url := range urls {
		var acg ACG
		parse, _ := url.Parse(_url)
		values, _ := url.ParseQuery(parse.RawQuery)
		// url query search number
		acg.SearchIndex = values.Get("s")

		dom, _ := getDecument(_url)
		s := dom.Find("div.ACG-mster_box1").First()
		acg.Image, _ = s.Find("img").Attr("src")
		acg.TaiName = s.Find("h1").First().Text()
		acg.JapName = s.Find("h2").First().Text()
		acg.Class = CheckColonExist(s.Find("ul.ACG-box1listA>li:contains(播映方式)").First().Text())
		acg.Premiere = CheckColonExist(s.Find("ul.ACG-box1listA>li:contains(當地首播)").First().Text())

		s.Find("ul.ACG-box1listB>li").Each(func(idx int, ss *goquery.Selection) {
			switch idx {
			case 0:
				// 爬取作者欄位
				acg.Author = CheckColonExist(ss.Text())
			case 1:
				// 爬取監督欄位
				acg.Director = CheckColonExist(ss.Text())
			case 2:
				// 爬取製作廠商
				acg.Firm = CheckColonExist(ss.Text())
			case 3:
				// 爬取台灣代理
				acg.Agent = CheckColonExist(ss.Text())
			case 4:
				// 爬取官方網站
				acg.Website = ss.Find("a").Text()
			}
		})
		log.Println(acg)
		time.Sleep(1 * time.Second)
		DB.Create(&acg)
	}
}

func CrewEachAnimeTest(_url string) {
	var acg ACG
	parse, _ := url.Parse(_url)
	values, _ := url.ParseQuery(parse.RawQuery)
	acg.SearchIndex = values.Get("s")

	dom, _ := getDecument(_url)
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

func GetAnimeInfo(_url string) (ACG, error) {
	var anime ACG
	parse, err := url.Parse(_url)
	if err != nil {
		return anime, err
	}

	values, err := url.ParseQuery(parse.RawQuery)
	if err != nil {
		return anime, err
	}

	searchidx := values.Get("s")

	DB.Where("search_index = ?", searchidx).First(&anime)
	return anime, nil
}

func SearchAnimeInfoWithindex(_url string) (ACG, error) {
	var anime ACG
	parse, _ := url.Parse(_url)

	values, _ := url.ParseQuery(parse.RawQuery)
	index := values.Get("s")
	err := DB.Where("search_index = ?", index).First(&anime).Error
	if err != nil {
		return anime, err
	}

	VerifyAnime(&anime)

	return anime, nil
}

func VerifyAnime(anime *ACG) {
	if len(anime.Agent) == 0 {
		anime.Agent = "nil"
	}
	if len(anime.Author) == 0 {
		anime.Author = "nil"
	}
	if len(anime.Class) == 0 {
		anime.Class = "nil"
	}
	if len(anime.Director) == 0 {
		anime.Director = "nil"
	}
	if len(anime.Firm) == 0 {
		anime.Firm = "nil"
	}
	if len(anime.Image) == 0 {
		anime.Image = "nil"
	}
	if len(anime.JapName) == 0 {
		anime.JapName = "nil"
	}
	if len(anime.Premiere) == 0 {
		anime.Premiere = "nil"
	}
	if len(anime.SearchIndex) == 0 {
		anime.SearchIndex = "nil"
	}
	if len(anime.TaiName) == 0 {
		anime.TaiName = "nil"
	}
	if len(anime.Website) == 0 {
		anime.Website = "nil"
	}
}

// 如果該作品有欄位為空, 填入nil以便於flex可以正常運作
func SearchAnimeInfoWithKey(key string) []ACG {
	var animes []ACG
	DB.Where("tai_name LIKE ?", "%"+key+"%").Limit(10).Find(&animes)
	for idx := 0; idx < len(animes); idx++ {
		VerifyAnime(&animes[idx])
	}
	return animes
}

func SearchAnimeInfoWithAuthor(author string) []ACG {
	var animes []ACG
	DB.Where("author LIKE ?", "%"+author+"%").Limit(10).Find(&animes)
	for idx := 0; idx < len(animes); idx++ {
		VerifyAnime(&animes[idx])
	}
	return animes
}
