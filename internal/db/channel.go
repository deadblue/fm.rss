package db

import (
	"context"
	"time"
)

type Channel struct {
	Source    string `bson:"source"`
	ChannelId string `bson:"channel_id"`

	Link   string `bson:"link"`
	Title  string `bson:"title"`
	Intro  string `bson:"intro"`
	Cover  string `bson:"cover"`
	Author string `bson:"author"`

	ItemCount  int  `bson:"item_count"`
	IsComplete bool `bson:"is_complete"`

	CreateTime time.Time `bson:"create_time"`
	UpdateTime time.Time `bson:"update_time"`
	FetchTime  time.Time `bson:"fetch_time"`
}

type _ChannelFilter struct {
	Source    string `bson:"source"`
	ChannelId string `bson:"channel_id"`
}

type _ChannelSetSpec struct {
	ItemCount  int       `bson:"item_count"`
	UpdateTime time.Time `bson:"update_time"`
	FetchTime  time.Time `bson:"fetch_time"`
}

func (c *Client) ChannelRetrieve(source, channelId string, channel *Channel) (err error) {
	result := c.channels.FindOne(context.Background(), &_ChannelFilter{
		Source:    source,
		ChannelId: channelId,
	})
	if err = result.Err(); err == nil {
		err = result.Decode(channel)
	}
	return
}

func (c *Client) ChannelInsert(doc *Channel) (err error) {
	_, err = c.channels.InsertOne(context.Background(), doc)
	return
}

func (c *Client) ChannelUpdate(doc *Channel) (err error) {
	_, err = c.channels.UpdateOne(context.Background(), &_ChannelFilter{
		Source:    doc.Source,
		ChannelId: doc.ChannelId,
	}, &_UpdateSpec{
		Set: &_ChannelSetSpec{
			ItemCount:  doc.ItemCount,
			UpdateTime: doc.UpdateTime,
			FetchTime:  doc.FetchTime,
		},
	})
	return
}
