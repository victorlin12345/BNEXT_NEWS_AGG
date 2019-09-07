package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"sync"
	"text/template"
	"time"
)

var INDEX_URL string = "https://www.bnext.com.tw/sitemap/google"

// Struct for Showing Data
type BNewsAggPage struct {
	PageTile string
	News     []Article
}

func newsRoutine(wg *sync.WaitGroup, c chan Articles, Location string) {
	defer wg.Done()
	var as Articles
	resp, _ := http.Get(Location)
	bytes, _ := ioutil.ReadAll(resp.Body)
	xml.Unmarshal(bytes, &as)
	resp.Body.Close()
	fmt.Printf("%s -> code status:%s amount:%d\n", Location, resp.Status, len(as.Titles))
	c <- as
}

func batch_process(c chan Articles, batch_data *[]string) {
	var wg sync.WaitGroup
	for _, data := range *batch_data {
		wg.Add(1)
		go newsRoutine(&wg, c, data)
	}
	wg.Wait()
	*batch_data = nil
}

func bnewsAggHandler(w http.ResponseWriter, r *http.Request) {

	var smi SitemapIndex
	var sitemap_count int = 50
	var batch_size int = 6 // Depending on the maximum request at the same time, Avoid code status:429
	var batch_data []string
	var queue chan Articles = make(chan Articles, 500)
	var showlist []Article
	var smooth_time time.Duration = 10

	start := time.Now()

	smi.FeedData(INDEX_URL)

	for iter, Loc := range smi.GetLocations(sitemap_count) {
		batch_data = append(batch_data, Loc)
		if iter != 0 && iter%(batch_size-1) == 0 {
			batch_process(queue, &batch_data)
			time.Sleep(time.Millisecond * smooth_time)
		}
		if iter == (sitemap_count - 1) {
			batch_process(queue, &batch_data)
		}
		//------- Sequiential Process --------
		/*
			var as Articles
			resp, _ := http.Get(Loc)
			bytes, _ := ioutil.ReadAll(resp.Body)
			xml.Unmarshal(bytes, &as)
			resp.Body.Close()
			fmt.Printf("%s -> code status:%s amount:%d\n", Loc, resp.Status, len(as.Titles))
			for idx, _ := range as.Titles {
				a := Article{Title: as.Titles[idx],
					Location:        as.Locations[idx],
					PucbicationDate: as.PucbicationDates[idx]}
				showlist = append(showlist, a)
			}
		*/
	}
	close(queue)

	for as := range queue {
		for idx, _ := range as.Titles {
			a := Article{Title: as.Titles[idx],
				Location:        as.Locations[idx],
				PucbicationDate: as.PucbicationDates[idx]}
			showlist = append(showlist, a)
		}
	}

	// sorted by PucbicationDate Lastest to Oldest
	sort.SliceStable(showlist, func(i, j int) bool {
		return showlist[i].PucbicationDate > showlist[j].PucbicationDate
	})

	elapsed := time.Since(start)
	log.Printf("Collecting %d data took %s\n", len(showlist), elapsed)

	// show on template
	p := BNewsAggPage{PageTile: "數位雜誌文章", News: showlist}
	t, _ := template.ParseFiles("agg.html")
	t.Execute(w, p)

}

func main() {

	http.HandleFunc("/", bnewsAggHandler)
	http.ListenAndServe(":8000", nil)

}
