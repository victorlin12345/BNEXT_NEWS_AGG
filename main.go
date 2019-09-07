package main

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"sort"
	"sync"
	"text/template"
)

var INDEX_URL string = "https://www.bnext.com.tw/sitemap/google"

// Struct for Showing Data
type BNewsAggPage struct {
	PageTile string
	News     []Article
}

var wg sync.WaitGroup

func newsRoutine(c chan Articles, Location string) {
	defer wg.Done()
	var as Articles
	resp, _ := http.Get(Location)
	bytes, _ := ioutil.ReadAll(resp.Body)
	xml.Unmarshal(bytes, &as)
	resp.Body.Close()
	c <- as
}

func bnewsAggHandler(w http.ResponseWriter, r *http.Request) {

	var smi SitemapIndex
	smi.FeedData(INDEX_URL)

	var showlist []Article
	var queue chan Articles = make(chan Articles, 100)

	for _, Loc := range smi.GetLocations(10) {
		wg.Add(1)
		go newsRoutine(queue, Loc)

		// var as Articles
		// resp, _ := http.Get(Loc)
		// bytes, _ := ioutil.ReadAll(resp.Body)
		// xml.Unmarshal(bytes, &as)
		// resp.Body.Close()

		// for idx, _ := range as.Titles {
		// 	a := Article{Title: as.Titles[idx],
		// 		Location:        as.Locations[idx],
		// 		PucbicationDate: as.PucbicationDates[idx]}
		// 	showlist = append(showlist, a)
		// }
	}

	wg.Wait()
	close(queue)

	for as := range queue {
		for idx, _ := range as.Titles {
			a := Article{Title: as.Titles[idx],
				Location:        as.Locations[idx],
				PucbicationDate: as.PucbicationDates[idx]}
			showlist = append(showlist, a)
		}
	}

	sort.SliceStable(showlist, func(i, j int) bool {
		return showlist[i].PucbicationDate > showlist[j].PucbicationDate
	})

	p := BNewsAggPage{PageTile: "數位雜誌文章", News: showlist}
	t, _ := template.ParseFiles("agg.html")
	t.Execute(w, p)

}

func main() {

	http.HandleFunc("/", bnewsAggHandler)
	http.ListenAndServe(":8000", nil)

}
