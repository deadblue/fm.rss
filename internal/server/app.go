package server

import "net/http"

type App interface {
	http.Handler
	Startup() (err error)
	Shutdown() (err error)
}
