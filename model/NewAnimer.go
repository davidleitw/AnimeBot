package model

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// https://acg.gamer.com.tw/quarterly.php?page=2&d=0
const NewAnimeUrl = "https://acg.gamer.com.tw/quarterly.php?"

type NewAnime struct {
	SearchIndex string `gorm:"primary_key; size:50; not null;"` // 動漫編號
	Popularity  int    // 人氣
	Followers   int    // 追隨數
	ImageSrc    string `gorm:"size:150;"` //首頁圖片
	TaiName     string `gorm:"size:90;"`  // 中文譯名
	JapName     string `gorm:"size:90;"`  // 原作名稱
	Class       string `gorm:"size:60;"`  // 動畫種類(電影or番)
	Premiere    string `gorm:"size:60;"`  // 首播時間
	Author      string `gorm:"size:60;"`  // 原著作者
	Director    string `gorm:"size:60;"`  // 導演監督
	Firm        string `gorm:"size:60;"`  // 製作廠商
	Agent       string `gorm:"size:60;"`  // 台灣代理
	Website     string `gorm:"size:150;"` // 官方網站
}

type NewAnimes []NewAnime

// Sort interface 三必要條件
func (na NewAnimes) Len() int { return len(na) }

func (na NewAnimes) Less(i, j int) bool {
	return (na[i].Popularity + na[i].Followers) > (na[j].Popularity + na[j].Followers)
}

func (na NewAnimes) Swap(i, j int) {
	na[i], na[j] = na[j], na[i]
}

func NewAnimeSortTest() {
	dbname := fmt.Sprintf("host=%s user=%s dbname=%s  password=%s", os.Getenv("HOST"), os.Getenv("DBUSER"), os.Getenv("DBNAME"), os.Getenv("PASSWORD"))
	ConnectDataBase(dbname)
	var nas NewAnimes
	DB.Find(&nas)
	fmt.Println(nas)
	sort.Sort(nas)
	for _, anime := range nas {
		fmt.Println(anime.Popularity)
	}
}

func CreateNewAnimeTable() {
	dbname := fmt.Sprintf("host=%s user=%s dbname=%s  password=%s", os.Getenv("HOST"), os.Getenv("DBUSER"), os.Getenv("DBNAME"), os.Getenv("PASSWORD"))
	ConnectDataBase(dbname)
	if DB.HasTable(&NewAnime{}) {
		DB.DropTable("new_animes")
	}
	DB.CreateTable(&NewAnime{})
	log.Println("Create new_animes table successfully!")
}

func CrewNewAnimeInfo() {
	dbname := fmt.Sprintf("host=%s user=%s dbname=%s  password=%s", os.Getenv("HOST"), os.Getenv("DBUSER"), os.Getenv("DBNAME"), os.Getenv("PASSWORD"))
	ConnectDataBase(dbname)

	newAnimePageNum := 2
	for i := 1; i <= newAnimePageNum; i++ {
		url := fmt.Sprintf("https://acg.gamer.com.tw/quarterly.php?page=%d&d=0", i)
		urls, animes_popularity := CrewNewAnimePageUrl(url)
		animes := CrewEachNewAnime(urls)
		for index, anime := range animes {
			anime.Popularity = animes_popularity[index]
			log.Println(anime)
			DB.Create(&anime)
		}
	}
}

func CrewNewAnimePageUrl(url string) ([]string, []int) {
	var urls []string
	var popularity []int

	dom, _ := getDecument(url)
	dom.Find("div.ACG-mainbox1").Each(func(idx int, s *goquery.Selection) {
		u, _ := s.Find("div.ACG-mainbox2>h1.ACG-maintitle>a").First().Attr("href")
		p, _ := strconv.Atoi(s.Find("div.ACG-mainbox4>p.ACG-mainplay>span").First().Text())

		urls = append(urls, "https:"+u)
		popularity = append(popularity, p)
	})
	return urls, popularity
}

func CrewEachNewAnime(urls []string) []NewAnime {
	var animes []NewAnime
	for _, _url := range urls {
		var anime NewAnime
		parse, _ := url.Parse(_url)
		values, _ := url.ParseQuery(parse.RawQuery)
		anime.SearchIndex = values.Get("s")

		dom, _ := getDecument(_url)
		anime.Followers, _ = strconv.Atoi(dom.Find("p.BH-acgbox>span").Text())
		s := dom.Find("div.ACG-mster_box1").First()
		anime.ImageSrc, _ = s.Find("img").Attr("src")
		anime.TaiName = CheckColonExist(s.Find("h1").First().Text())
		anime.JapName = CheckColonExist(s.Find("h2").First().Text())
		if len(anime.JapName) == 0 {
			anime.JapName = "nil"
		}
		anime.Class = CheckColonExist(s.Find("ul.ACG-box1listA>li:contains(播映方式)").First().Text())
		anime.Premiere = CheckColonExist(s.Find("ul.ACG-box1listA>li:contains(當地首播)").First().Text())

		s.Find("ul.ACG-box1listB>li").Each(func(idx int, ss *goquery.Selection) {
			switch idx {
			case 0:
				// 爬取作者欄位
				anime.Author = CheckColonExist(ss.Text())
				if len(anime.Author) == 0 {
					anime.Author = "nil"
				}
			case 1:
				// 爬取監督欄位
				anime.Director = CheckColonExist(ss.Text())
				if len(anime.Director) == 0 {
					anime.Director = "nil"
				}
			case 2:
				// 爬取製作廠商
				anime.Firm = CheckColonExist(ss.Text())
				if len(anime.Firm) == 0 {
					anime.Firm = "nil"
				}
			case 3:
				// 爬取台灣代理
				anime.Agent = CheckColonExist(ss.Text())
				if len(anime.Agent) == 0 {
					anime.Agent = "nil"
				}
			case 4:
				// 爬取官方網站
				anime.Website = ss.Find("a").Text()
				if len(anime.Website) == 0 {
					anime.Website = "nil"
				}
			}
		})
		time.Sleep(1 * time.Second)
		animes = append(animes, anime)
	}
	return animes
}
