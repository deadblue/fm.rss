package app

import (
	"errors"
	"flag"
	"os"
)

const (
	_EnvMongoUri   = "MONGO_URI"
	_EnvServerNet  = "SERVER_NET"
	_EnvServerAddr = "SERVER_ADDR"

	_FlagMongoUri   = "mongo-uri"
	_FlagServerNet  = "server-net"
	_FlagServerAddr = "server-addr"

	_DefaultServerNet  = "tcp"
	_DefaultServerAddr = ":9066"
)

var (
	errInvalidOptions = errors.New("invalid options")
)

type _Options struct {
	MongoUri   string
	ServerNet  string
	ServerAddr string
}

func (o *_Options) AutoConf() *_Options {
	// Get from environment variable first
	defaultMongoUri := os.Getenv(_EnvMongoUri)
	defaultServerNet := os.Getenv(_EnvServerNet)
	if defaultServerNet == "" {
		defaultServerNet = _DefaultServerNet
	}
	defaultServerAddr := os.Getenv(_EnvServerAddr)
	if defaultServerAddr == "" {
		defaultServerAddr = _DefaultServerAddr
	}
	// Get from command line
	flag.StringVar(&o.MongoUri, _FlagMongoUri, defaultMongoUri, "Mongo URI")
	flag.StringVar(&o.ServerNet, _FlagServerNet, defaultServerNet, "Server network")
	flag.StringVar(&o.ServerAddr, _FlagServerAddr, defaultServerAddr, "Server address")
	flag.Parse()
	return o
}

func (o *_Options) SelfCheck() error {
	if o.MongoUri == "" {
		return errInvalidOptions
	}
	return nil
}
