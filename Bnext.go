package main

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

type SitemapIndex struct {
	Locations []string `xml:"sitemap>loc"`
	LastMods  []string `xml:"sitemap>lastmod"`
}

// Methods - Pointer Receivers
func (smi *SitemapIndex) FeedData(index_url string) {
	resp, _ := http.Get(index_url)
	bytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	xml.Unmarshal(bytes, &smi)
}

// Methods - Pointer Receivers
func (smi *SitemapIndex) GetLocations(max_len int) []string {
	var LocationList []string
	for idx, Loc := range smi.Locations {
		if idx != 0 && idx <= max_len {
			LocationList = append(LocationList, Loc)
		}
	}
	return LocationList
}

// struct for Articles
type Articles struct {
	Locations        []string `xml:"url>loc"`
	PucbicationDates []string `xml:"url>news>publication_date"`
	Titles           []string `xml:"url>news>title"`
}

type Article struct {
	Title           string
	Location        string
	PucbicationDate string
}
