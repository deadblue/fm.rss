package db

import (
	"context"
	"log"
	"time"
)

type Item struct {
	Source    string `bson:"source"`
	ItemId    string `bson:"item_id"`
	ChannelId string `bson:"channel_id"`

	Guid  string `bson:"guid"`
	Link  string `bson:"link"`
	Title string `bson:"title"`
	Intro string `bson:"intro"`

	Url      string `bson:"url"`
	MimeType string `bson:"mime_type"`
	Size     int64  `bson:"size"`
	Duration int    `bson:"duration"`

	CreateTime time.Time `bson:"create_time"`
	FetchTime  time.Time `bson:"fetch_time"`
}

type _ItemFilter struct {
	Source    string `bson:"source"`
	ChannelId string `bson:"channel_id"`
}

func (c *Client) ItemRetrieve(source, channelId string) (items []*Item, err error) {
	cursor, err := c.items.Find(context.Background(), &_ItemFilter{
		Source:    source,
		ChannelId: channelId,
	})
	if err != nil {
		return
	}
	defer func() {
		_ = cursor.Close(context.TODO())
	}()
	items = make([]*Item, 0)
	for cursor.Next(context.TODO()) {
		item := &Item{}
		if errDec := cursor.Decode(item); errDec == nil {
			items = append(items, item)
		} else {
			log.Printf("Decode document error: %s", errDec)
		}
	}
	return
}

func (c *Client) ItemInsert(items []*Item) (err error) {
	// Convert into interface{} slice.
	docs := make([]interface{}, len(items))
	for index, item := range items {
		docs[index] = item
	}
	// Save to database
	_, err = c.items.InsertMany(context.Background(), docs)
	return
}
