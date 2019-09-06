package main

import "fmt"

var INDEX_URL string = "https://www.bnext.com.tw/sitemap/google"

func main() {

	var smi SitemapIndex
	smi.FeedData(INDEX_URL)
	for _, Loc := range smi.GetLocations() {
		fmt.Println(Loc)
	}

}
