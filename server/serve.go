package server

import (
	"context"
	"fmt"
	"log"
	"main/job/metrics"
	"main/worker"
	"net"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func Serve(ctx context.Context, maxWorkTime time.Duration, port uint16, reg *prometheus.Registry, metrics *metrics.Metrics, worker *worker.BasicWorker) error {
	var srvCtx context.Context
	var cancel context.CancelFunc

	if maxWorkTime > 0 {
		srvCtx, cancel = context.WithTimeout(ctx, maxWorkTime)
	} else {
		srvCtx, cancel = context.WithCancel(ctx)
	}

	defer cancel()

	mux := SetAndGetMux(reg, metrics, worker)
	addr := fmt.Sprintf("localhost:%d", port)
	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			return srvCtx
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
	case <-srvCtx.Done():
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
