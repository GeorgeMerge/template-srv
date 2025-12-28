package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	http.Server

	m    sync.Mutex
	port int
}

func New(ctx context.Context, router http.Handler, port int) *Server {
	const defaultReadTimeout = time.Second * 30

	return &Server{
		Server: http.Server{
			Addr:    fmt.Sprintf(":%v", port),
			Handler: router,
			BaseContext: func(_ net.Listener) context.Context {
				return ctx
			},
			ReadTimeout: defaultReadTimeout,
		},
		port: port,
	}
}

func (s *Server) Port() int {
	s.m.Lock()
	defer s.m.Unlock()

	return s.port
}

func (s *Server) ListenAndServe() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", s.port))
	if err != nil {
		return err
	}

	if tcpAddr, ok := listener.Addr().(*net.TCPAddr); ok {
		s.setPort(tcpAddr.Port)
	} else {
		return errors.New("error getting tcp address")
	}

	return s.Serve(listener)
}

func (s *Server) setPort(port int) {
	s.m.Lock()
	defer s.m.Unlock()
	s.port = port
}
