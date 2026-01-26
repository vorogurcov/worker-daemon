package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func setAndGetMux(reg *prometheus.Registry) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))

	return mux
}

func Serve(ctx context.Context, port uint16, reg *prometheus.Registry) error {

	mux := setAndGetMux(reg)
	addr := fmt.Sprintf("localhost:%d", port)
	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			return ctx
		},
	}

	errCh := make(chan error, 1)

	go func() {
		log.Printf("http: starting server on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		log.Printf("http: shutting down server (graceful)...")
		if err := srv.Shutdown(shutdownCtx); err != nil {
			_ = srv.Close()
			return fmt.Errorf("server shutdown failed: %w", err)
		}
		return nil
	case err := <-errCh:
		if err != nil {
			return err
		}
		return nil
	}
}
