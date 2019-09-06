package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

var INDEX_URL string = "https://www.bnext.com.tw/sitemap/google"

func main() {

	var smi SitemapIndex
	smi.FeedData(INDEX_URL)
	for _, Loc := range smi.GetLocations() {
		var articles Articles
		resp, _ := http.Get(Loc)
		bytes, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		xml.Unmarshal(bytes, &articles)
		fmt.Println(articles.Titles)
	}

}
