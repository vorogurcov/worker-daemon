package server

import (
	"context"
	"fmt"
	"net/http"
)

func setAndGetMux() *http.ServeMux {
	// Если не ставить * перед http.ServeMux, то оно скопирует его?
	//TODO: Разобраться, куда тут пихать контекст. Расширять r.Context() ?
	mux := http.NewServeMux()

	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {

	})

	return mux
}

func Serve(ctx context.Context, port uint8) error {
	mux := setAndGetMux()

	if err := http.ListenAndServe(fmt.Sprintf("localhost:%v", port), mux); err != nil {
		return err
	}

	return nil
}
