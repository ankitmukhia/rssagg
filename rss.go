// RSS is just structured data in XML format, Kind of like crappy JSON.
package main

import (
	"net/http"
	"time"
	"encoding/xml"
	"io"
)

type RSSFeed struct {
	Channel struct {
		Title string `xml:"title"`
		Link string `xml:"link"`
		Description string `xml:"description"`
		Language string `xml:"language"`
		Item []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title string `xml:"title"`	
	Link string `xml:"link"`
	Description string `xml:"description"`
	PubDate string `xml:"pubDate"` 
}

func urlToFeed(url string) (RSSFeed, error) {
	// fetch every ten second
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}

	res, err := httpClient.Get(url)
	if err != nil {
		return RSSFeed{}, nil
	}
	defer res.Body.Close()
	
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return RSSFeed{}, nil
	}

	rssFeed := RSSFeed{}

	err = xml.Unmarshal(data, &rssFeed)
	if err != nil {
		return RSSFeed{}, nil
	}
	
	return rssFeed, nil
}
