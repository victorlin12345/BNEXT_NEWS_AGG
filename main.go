package main

import "fmt"

var INDEX_URL string = "https://www.bnext.com.tw/sitemap/google"

func main() {

	var smi SitemapIndex
	smi.FeedData(INDEX_URL)

	for idx, Loc := range smi.Locations {
		if idx != 0 {
			fmt.Println(Loc)
		}
	}
}
