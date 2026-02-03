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

func Serve(
	ctx context.Context,
	maxWorkTime time.Duration,
	port uint16,
	reg *prometheus.Registry,
	metrics *metrics.Metrics,
	worker *worker.BasicWorker,
) error {

	var srvCtx context.Context
	var cancel context.CancelFunc

	if maxWorkTime > 0 {
		srvCtx, cancel = context.WithTimeout(ctx, maxWorkTime)
	} else {
		srvCtx, cancel = context.WithCancel(ctx)
	}
	defer cancel()

	mux := SetAndGetMux(reg, metrics, worker)

	srv := &http.Server{
		Addr:    fmt.Sprintf("localhost:%d", port),
		Handler: mux,
		BaseContext: func(net.Listener) context.Context {
			return srvCtx
		},
	}

	errCh := make(chan error, 1)

	go func() {
		log.Printf("http: starting on %s", srv.Addr)
		errCh <- srv.ListenAndServe()
	}()

	select {
	case <-srvCtx.Done():
		log.Println("http: graceful shutdown...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			_ = srv.Close()
			return fmt.Errorf("shutdown failed: %w", err)
		}

		err := <-errCh
		if err != nil && err != http.ErrServerClosed {
			return err
		}

		log.Println("http: shutdown complete")
		return nil

	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	}
}
