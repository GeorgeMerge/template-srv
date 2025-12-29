package httpserver

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

type Server struct {
	http.Server
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

func (s *Server) ListenAndServe() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", s.port))
	if err != nil {
		return err
	}

	return s.Serve(listener)
}
