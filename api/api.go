package api

import (
	"context"
	"d7024e_group04/env"
	"d7024e_group04/internal/node"
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
	mux := http.NewServeMux()

	mux.HandleFunc("GET /object/", h.getObject())
	mux.HandleFunc("POST /object", h.putObject())

	return http.ListenAndServe(":"+strconv.Itoa(env.ApiPort), mux)
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
