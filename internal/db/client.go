package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

const (
	_DbName = "podcast"

	_CollItem    = "item"
	_CollChannel = "channel"
)

type Client struct {
	mc *mongo.Client
	db *mongo.Database

	items    *mongo.Collection
	channels *mongo.Collection
}

func (c *Client) Init() (err error) {
	log.Println("Connecting to mongo DB ...")
	if err = c.mc.Connect(context.Background()); err == nil {
		// Get collection handles
		c.db = c.mc.Database(_DbName)
		c.items = c.db.Collection(_CollItem)
		c.channels = c.db.Collection(_CollChannel)
	}
	return
}

func (c *Client) Release() error {
	log.Println("Disconnecting from mongo DB ...")
	return c.mc.Disconnect(context.Background())
}

func New(uri string) (client *Client, err error) {
	mc, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err == nil {
		client = &Client{
			mc: mc,
		}
	}
	return
}
