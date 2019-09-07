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
	News     []Article
}

func bnewsAggHandler(w http.ResponseWriter, r *http.Request) {

	var smi SitemapIndex
	smi.FeedData(INDEX_URL)

	var showlist []Article

	for _, Loc := range smi.GetLocations(5) {

		var articles Articles
		resp, _ := http.Get(Loc)
		bytes, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		xml.Unmarshal(bytes, &articles)

		for idx, _ := range articles.Locations {

			a := Article{Title: articles.Titles[idx],
				Location:        articles.Locations[idx],
				PucbicationDate: articles.PucbicationDates[idx]}

			showlist = append(showlist, a)
		}
	}

	p := BNewsAggPage{PageTile: "數位雜誌文章", News: showlist}
	t, _ := template.ParseFiles("agg.html")
	t.Execute(w, p)

}

func main() {

	http.HandleFunc("/", bnewsAggHandler)
	http.ListenAndServe(":8000", nil)

}
