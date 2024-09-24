package api

import (
	"context"
	"d7024e_group04/env"
	"d7024e_group04/internal/node"
	"log"
	"net/http"
	"strconv"
)

type Handler struct {
	node *node.Node
}

func NewHandler(node *node.Node) *Handler {
	return &Handler{
		node: node,
	}
}

func (h *Handler) ListenAndServe(ctx context.Context) error {
	errChan := make(chan error, 1)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /object/", h.getObject())
	mux.HandleFunc("POST /object", h.putObject())

	go func() {
		log.Printf("serving http requests on port %v", env.ApiPort)
		errChan <- http.ListenAndServe(":"+strconv.Itoa(env.ApiPort), mux)
	}()

	for {
		select {
		case err := <-errChan:
			return err
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (h *Handler) getObject() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		panic("TODO")
	}
}

func (h *Handler) putObject() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		panic("TODO")
	}
}
