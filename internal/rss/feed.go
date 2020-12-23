package rss

import (
	"bytes"
	"encoding/xml"
	"io"
	"strconv"
	"time"
)

const (
	_DateLayout = time.RFC1123
)

func (f *Feed) WithInfo(link, title, description string) *Feed {
	// Standard properties
	f.Channel.Link = link
	f.Channel.Title = title
	f.Channel.Description = description
	if f.Channel.Image != nil {
		f.Channel.Image.Title = title
		f.Channel.Image.Link = link
	}
	// iTunes properties
	f.Channel.ItunesTitle = title
	return f
}

func (f *Feed) WithAuthor(author string) *Feed {
	f.Channel.ItunesAuthor = author
	return f
}

func (f *Feed) WithImage(image string) *Feed {
	// Standard properties
	f.Channel.Image = &_Image{
		Url:   image,
		Title: f.Channel.Title,
		Link:  f.Channel.Link,
	}
	// iTunes properties
	f.Channel.ItunesImage.Href = image
	return f
}

func (f *Feed) WithComplete(complete bool) *Feed {
	if complete {
		f.Channel.ItunesComplete = "Yes"
	}
	return f
}

func (f *Feed) WithPubDate(pubDate time.Time) *Feed {
	f.Channel.PubDate = pubDate.Format(_DateLayout)
	return f
}

func (f *Feed) NewItem() *Item {
	item := &Item{}
	f.Channel.Item = append(f.Channel.Item, item)
	return item
}

func (f *Feed) WriteTo(w io.Writer) (n int64, err error) {
	data, err := xml.Marshal(f)
	if err != nil {
		return
	}
	nn, err := w.Write([]byte(xml.Header))
	if err != nil {
		return
	} else {
		n += int64(nn)
	}
	nn, err = w.Write(data)
	if err != nil {
		return
	} else {
		n += int64(nn)
	}
	return
}

func (f *Feed) String() string {
	// Write header
	buf := &bytes.Buffer{}
	buf.WriteString(xml.Header)
	// Write body
	body, _ := xml.Marshal(f)
	_, _ = buf.Write(body)
	return buf.String()
}

func (i *Item) WithInfo(guid, link, title, description string) *Item {
	i.Guid = guid
	i.Link = link
	i.Title = title
	i.Description = description

	i.ItunesTitle = title
	return i
}

func (i *Item) WithEnclosure(url string, size int64, mimeType string) *Item {
	i.Enclosure.Url = url
	i.Enclosure.Type = mimeType
	i.Enclosure.Length = strconv.FormatInt(size, 10)
	return i
}

func (i *Item) WithDuration(seconds int) *Item {
	i.ItunesDuration = strconv.Itoa(seconds)
	return i
}

func (i *Item) WithPubDate(pubDate time.Time) *Item {
	i.PubDate = pubDate.Format(_DateLayout)
	return i
}

func New() *Feed {
	return &Feed{
		XMLName:  xml.Name{},
		NsItunes: "http://www.itunes.com/dtds/podcast-1.0.dtd",
		Version:  "2.0",
		Channel: &_Channel{
			Item: make([]*Item, 0),
		},
	}
}
