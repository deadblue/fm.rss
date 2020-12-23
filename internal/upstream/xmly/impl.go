package xmly

import (
	"github.com/deadblue/fm.rss/internal/upstream"
	"net/http"
)

const (
	_Name = "ximalaya"
)

type fetcherImpl struct {
	hc *http.Client
}

func (f *fetcherImpl) Name() string {
	return _Name
}

func (f *fetcherImpl) Aliases() []string {
	return []string{
		"xmly", "ximalaya",
	}
}

func New() upstream.Fetcher {
	return &fetcherImpl{
		hc: &http.Client{},
	}
}
