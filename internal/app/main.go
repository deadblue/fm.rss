package app

import (
	"github.com/deadblue/fm.rss/internal/server"
	"log"
)

func Main() {
	// Fill options
	opts := (&_Options{}).AutoConf()
	if err := opts.SelfCheck(); err != nil {
		log.Fatal(err)
	}
	// Create app
	a, err := create(opts.MongoUri)
	if err != nil {
		log.Fatal(err)
	}
	// Create server and run
	s, err := server.New(opts.ServerNet, opts.ServerAddr, a)
	if err != nil {
		log.Fatal(err)
	}
	s.Run()
}
