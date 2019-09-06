package main

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"text/template"
)

var INDEX_URL string = "https://www.bnext.com.tw/sitemap/google"

// Struct for Showing Data
type BNewsAggPage struct {
	PageTile string
	News     map[string]Article
}

func bnewsAggHandler(w http.ResponseWriter, r *http.Request) {

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

	p := BNewsAggPage{PageTile: "數位雜誌文章", News: article_map}
	t, _ := template.ParseFiles("agg.html")
	t.Execute(w, p)

}

func main() {

	http.HandleFunc("/", bnewsAggHandler)
	http.ListenAndServe(":8000", nil)

}
