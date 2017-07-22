package main

import (
	"github.com/mmcdole/gofeed"
	"fmt"
	"strconv"

)

// rss1ページ辺りの件数は固定
const perPage = 20
// timeformat

func main() {
	//rss feed
	rssurl := "http://b.hatena.ne.jp/ytacky/rss"
	// indexは1から開始
	index := 1
	id := 0
	fp := gofeed.NewParser()
	for {
		BookmarkFeed, _ := fp.ParseURL(rssurl + "?of=" + strconv.Itoa(index))
		index = index + perPage
		items := BookmarkFeed.Items
		if len(items) == 0 {
			fmt.Println("End Of purse")
			break
		}
		for _, item := range items {
			fmt.Println("bookmarkId:" + strconv.Itoa(id))
			fmt.Println("Title: " + item.Title)
			fmt.Println("Link:" + item.Link)
			categories := item.Categories
			category_num := 0
			for _, category := range categories {
				fmt.Println()
				fmt.Println("Category_id : " + strconv.Itoa(category_num))
				fmt.Println("category : " + category)
				category_num+
			}

			bookmarkDate := item.Extensions["dc"]["date"][0].Value
			fmt.Println(bookmarkDate)

			break
		}
	}

}
