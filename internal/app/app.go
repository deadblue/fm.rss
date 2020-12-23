package app

import (
	"github.com/deadblue/fm.rss/internal/db"
	"github.com/deadblue/fm.rss/internal/server"
	"github.com/deadblue/fm.rss/internal/upstream"
	"github.com/deadblue/fm.rss/internal/upstream/xmly"
	"log"
)

type appImpl struct {
	// Database client
	dc *db.Client
	// Upstream fetchers
	fs map[string]upstream.Fetcher
}

func (a *appImpl) Startup() (err error) {
	return a.dc.Init()
}

func (a *appImpl) Shutdown() (err error) {
	return a.dc.Release()
}

func (a *appImpl) Register(fetchers ...upstream.Fetcher) {
	if a.fs == nil {
		a.fs = make(map[string]upstream.Fetcher)
	}
	for _, fetcher := range fetchers {
		log.Printf("Register upstream fetcher: %s", fetcher.Name())
		for _, alias := range fetcher.Aliases() {
			a.fs[alias] = fetcher
		}
	}
}

func create(dbUri string) (_ server.App, err error) {
	dc, err := db.New(dbUri)
	if err != nil {
		return
	}
	impl := &appImpl{
		dc: dc,
		fs: make(map[string]upstream.Fetcher),
	}
	// Register fetchers
	impl.Register(xmly.New())
	return impl, nil
}
