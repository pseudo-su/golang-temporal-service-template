package httpserver

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"time"

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
	go func() {
		if err := s.combinedServer.Serve(listener); err != nil && err != grpc.ErrServerStopped && err != http.ErrServerClosed {
			slog.ErrorContext(context.Background(), "shutdown error", slog.Any("error", err))
		}
	}()

	<-interruptCh

	gracefullCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if err := s.combinedServer.Shutdown(gracefullCtx); err != nil {
		return err
	}

	return nil
}
