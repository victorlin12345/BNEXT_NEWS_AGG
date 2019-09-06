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

	article_map := make(map[string]Article)

	for _, Loc := range smi.GetLocations() {
		var articles Articles
		resp, _ := http.Get(Loc)
		bytes, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		xml.Unmarshal(bytes, &articles)
		for idx, _ := range articles.Locations {
			a := Article{Location: articles.Locations[idx],
				PucbicationDate: articles.PucbicationDates[idx]}
			article_map[articles.Titles[idx]] = a
		}
	}

	for t, a := range article_map {
		fmt.Printf("\n\n\n%s", t)
		fmt.Printf("\n%s", a.Location)
		fmt.Printf("\n%s", a.PucbicationDate)
	}

}
