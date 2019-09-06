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
