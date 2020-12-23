package rss

import (
	"encoding/xml"
)

type Item struct {
	Guid        string `xml:"guid"`
	Link        string `xml:"link"`
	Title       string `xml:"title"`
	Description string `xml:"description"`

	Enclosure struct {
		Url    string `xml:"url,attr"`
		Length string `xml:"length,attr"`
		Type   string `xml:"type,attr"`
	} `xml:"enclosure"`

	PubDate string `xml:"pubDate"`

	ItunesTitle    string `xml:"itunes:title"`
	ItunesDuration string `xml:"itunes:duration"`
}

type _Image struct {
	Url   string `xml:"url"`
	Title string `xml:"title"`
	Link  string `xml:"link"`
}

type _Channel struct {
	Title       string  `xml:"title"`
	Link        string  `xml:"link"`
	Description string  `xml:"description"`
	Image       *_Image `xml:"image,omitempty"`
	PubDate     string  `xml:"pubDate"`

	ItunesTitle  string `xml:"itunes:title"`
	ItunesAuthor string `xml:"itunes:author"`
	ItunesImage  struct {
		Href string `xml:"href,attr"`
	} `xml:"itunes:image"`
	ItunesComplete string `xml:"itunes:complete,omitempty"`

	Item []*Item `xml:"item"`
}

type Feed struct {
	XMLName  xml.Name  `xml:"rss"`
	Version  string    `xml:"version,attr"`
	NsItunes string    `xml:"xmlns:itunes,attr"`
	Channel  *_Channel `xml:"channel"`
}
