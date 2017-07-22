package main

import (
	"github.com/mmcdole/gofeed"
	"fmt"
	"strconv"
	"os"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"

	"log"
	"github.com/joho/godotenv"
)

type bookmark struct {
	gorm.Model
	Id    int `gorm:"primary_key"`
	Title string
	Datetime string
	Link string
}

type tag struct {
	gorm.Model
	Tag_Id int `gorm:"primary_key"`
	Bookmark_Id int
	Tag string

}
// rss1ページ辺りの件数は固定
const perPage = 20

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	rssurl := os.Getenv("RSS_URL")

	// rssのstart indexは1
	index := 1

	// connect DB
	fp := gofeed.NewParser()
	DB_HOST := os.Getenv("DB_HOST")
	DB_CHARSET := os.Getenv("DB_CHARSET")
	DB_USER := os.Getenv("DB_USER")
	DB_PASS := os.Getenv("DB_PASS")
	DB_NAME := os.Getenv("DB_NAME")
	DB_PORT := os.Getenv("DB_PORT")
	DB_CONNECT := DB_USER + ":" + DB_PASS + "@tcp(" + DB_HOST + ":" + DB_PORT + ")/" + DB_NAME +"?charset=" + DB_CHARSET + "&parseTime=true&loc=Asia%2FTokyo"
	db, err := gorm.Open("mysql", DB_CONNECT)
	defer db.Close()

	if err != nil {
		log.Println(err)
		fmt.Println("Can't Connect DB")
		os.Exit(1)
	}

	for {
		bookmarksIn := bookmark{}
		BookmarkFeed, _ := fp.ParseURL(rssurl + "?of=" + strconv.Itoa(index))
		items := BookmarkFeed.Items

		// rssが0件ならば終了
		if len(items) == 0 {
			fmt.Println("End Of purse")
			os.Exit(0)
		}

		for _, item := range items {
			bookmarkDate := item.Extensions["dc"]["date"][0].Value
			// RFC3339形式なのでdatetimeで扱える形式に変換する
			t, _ := time.Parse(time.RFC3339, bookmarkDate)


			bookmarksIn.ID = 0
			bookmarksIn.Title = item.Title
			bookmarksIn.Link = item.Link
			bookmarksIn.Datetime = t.Format("2006-01-02 15:04:05")
			db.Create(&bookmarksIn)

			fmt.Println("bookmarkId:" + strconv.Itoa(int(bookmarksIn.ID)))
			fmt.Println("Title: " + item.Title)
			fmt.Println("Link:" + item.Link)
			fmt.Println("date: " + t.Format("2006-01-02 15:04:05"))

			categories := item.Categories
			for _, category := range categories {
				tagIn := tag{}
				tagIn.Bookmark_Id = int(bookmarksIn.ID)
				tagIn.Tag = category
				tagIn.Tag_Id = 0
				fmt.Println()
				fmt.Println("bookmarkId : " +strconv.Itoa(int(bookmarksIn.ID)))
				fmt.Println("Tag : " + category)
				db.Create(&tagIn)
				fmt.Println("tag unique" + strconv.Itoa(int(tagIn.ID)))
			}
		}

		// perPage分indexを進める
		index = index + perPage

	}

}

