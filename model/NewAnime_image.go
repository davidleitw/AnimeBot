package model

import (
	"fmt"
	"os"
)

type updateNewAnime struct {
	index string
	image string
	popu  int
}

var updateData []updateNewAnime = []updateNewAnime{
	{
		index: "112147",
		image: "https://p2.bahamut.com.tw/B/2KU/11/4b5ff0d2b86f56ac95ab79af7c1a2sz5.JPG",
		popu:  2018,
	},
	{
		index: "113501",
		image: "",
		popu:  1408,
	},
	{
		index: "113790",
		image: "https://p2.bahamut.com.tw/WIKI/14/00384714.JPG",
		popu:  839,
	},
	{
		index: "108065",
		image: "https://p2.bahamut.com.tw/WIKI/70/00384270.JPG",
		popu:  145,
	},
	{
		index: "109207",
		image: "https://p2.bahamut.com.tw/B/2KU/20/ed5f83770965c627d39fc522321a2w05.JPG",
		popu:  64,
	},
	{
		index: "110182",
		image: "https://p2.bahamut.com.tw/B/2KU/53/16030bf3d383c173a00eb9fc6419r155.JPG",
		popu:  54,
	},
	{
		index: "111051",
		image: "https://p2.bahamut.com.tw/WIKI/55/00384155.JPG",
		popu:  78,
	},
	{
		index: "111282",
		image: "https://p2.bahamut.com.tw/B/2KU/88/632b07cbf087ac82ed6b1df31419hss5.JPG",
		popu:  120,
	},
	{
		index: "113532",
		image: "https://p2.bahamut.com.tw/B/2KU/94/28ba4c69324195f92bc3cedb2719gq25.JPG",
		popu:  45,
	},
	{
		index: "109387",
		image: "https://p2.bahamut.com.tw/WIKI/86/00384286.JPG",
		popu:  1575,
	},
	{
		index: "108139",
		image: "https://p2.bahamut.com.tw/WIKI/38/00384138.JPG",
		popu:  64,
	},
	{
		index: "109451",
		image: "https://p2.bahamut.com.tw/B/2KU/22/4942a12be7df3d9ba97a2dfdaa19w0u5.JPG",
		popu:  8584,
	},
	{
		index: "107955",
		image: "https://p2.bahamut.com.tw/WIKI/27/00384327.JPG",
		popu:  6662,
	},
	{
		index: "108269",
		image: "https://p2.bahamut.com.tw/WIKI/02/00384202.JPG",
		popu:  7478,
	},
	{
		index: "108886",
		image: "https://p2.bahamut.com.tw/B/2KU/04/cb4ac714360f311d328262fa7d19w345.JPG",
		popu:  7500,
	},
	{
		index: "109015",
		image: "https://p2.bahamut.com.tw/B/2KU/58/6183e5ae9e495289015e51c58218mde5.JPG",
		popu:  821,
	},
	{
		index: "110208",
		image: "",
		popu:  5939,
	},
	{
		index: "110246",
		image: "https://p2.bahamut.com.tw/B/2KU/37/f7512c50b2f7946d520f26309f19w415.JPG",
		popu:  91,
	},
	{
		index: "110536",
		image: "https://p2.bahamut.com.tw/B/2KU/79/0665ed837ac66c1ba96bbf6a7b19mt35.JPG",
		popu:  6953,
	},
}

func UpdateNewAnimeImage() {
	dbname := fmt.Sprintf("host=%s user=%s dbname=%s  password=%s", os.Getenv("HOST"), os.Getenv("DBUSER"), os.Getenv("DBNAME"), os.Getenv("PASSWORD"))
	ConnectDataBase(dbname)

	for _, data := range updateData {
		if data.image != "" {
			DB.Where("search_index = ?", data.index).Update("image_src", data.image)
		}
		if data.popu != 0 {
			DB.Where("search_index = ?", data.index).Update("popularity", data.popu)
		}
	}
}
