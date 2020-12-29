package server

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
)

type Server struct {
	nl  net.Listener
	hs  *http.Server
	app App

	cf int32
	ch chan error
}

func (s *Server) close(err error) {
	// Avoid double closing
	if atomic.CompareAndSwapInt32(&s.cf, 0, 1) {
		if err != nil {
			s.ch <- err
		}
		close(s.ch)
	}
}

func (s *Server) startup() {
	go func() {
		// Start app before server
		if err := s.app.Startup(); err != nil {
			s.close(err)
			return
		}
		// Start server
		addr := s.nl.Addr()
		log.Printf("HTTP server listening at: %s://%s", addr.Network(), addr.String())
		if err := s.hs.Serve(s.nl); err != http.ErrServerClosed {
			s.close(err)
		}
	}()
}

func (s *Server) shutdown() {
	go func() {
		if atomic.LoadInt32(&s.cf) == 0 {
			// Shutting down server
			_ = s.hs.Shutdown(context.Background())
			// Shutting down app
			err := s.app.Shutdown()
			s.close(err)
		}
	}()
}

func (s *Server) Run() {
	// Handle OS signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	// Startup server
	s.startup()
	// Waiting for down
	for alive := true; alive; {
		select {
		case <-sigCh:
			log.Printf("Server shutting down for signal ...")
			s.shutdown()
		case err, isRecv := <-s.ch:
			if isRecv {
				log.Printf("Server shutting down due to error: %s ...", err)
			} else {
				log.Printf("Server shutdown!")
			}
			alive = isRecv
		}
	}
}

func New(network, address string, app App) (s *Server, err error) {
	if network == "unix" || network == "unixpacket" {
		if err = os.Remove(address); err != nil && !os.IsNotExist(err) {
			return
		}
	}
	// Create network listener
	l, err := net.Listen(network, address)
	if err != nil {
		return
	}
	if ul, ok := l.(*net.UnixListener); ok {
		ul.SetUnlinkOnClose(true)
	}
	// Create server
	s = &Server{
		nl: l,
		hs: &http.Server{
			Handler: app,
		},
		app: app,

		cf: 0,
		ch: make(chan error),
	}
	return
}
