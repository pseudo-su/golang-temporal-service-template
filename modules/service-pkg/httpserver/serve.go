package httpserver

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

func (s *HttpServer) ListenAndServe(interruptCh <-chan interface{}) error {
	var listener net.Listener = s.opts.listener
	if listener == nil {
		l, err := net.Listen("tcp", s.opts.address)
		if err != nil {
			return err
		}
		listener = l
	}

	return s.Serve(listener, interruptCh)
}

func (s *HttpServer) Serve(listener net.Listener, interruptCh <-chan interface{}) error {
	// https://github.com/golang/go/issues/26682#issuecomment-431904745
	h2s := &http2.Server{}
	h1s := &http.Server{}
	err := http2.ConfigureServer(h1s, h2s)
	if err != nil {
		return err
	}
	go func() {
		if err := http.Serve(listener, h2c.NewHandler(s.connectServer, h2s)); err != nil && err != grpc.ErrServerStopped && err != http.ErrServerClosed {
			slog.ErrorContext(context.Background(), "shutdown error", slog.Any("error", err))
		}
	}()

	<-interruptCh

	gracefullCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if err := h1s.Shutdown(gracefullCtx); err != nil {
		return err
	}

	return nil
}
