package api

import (
	"context"
	"d7024e_group04/env"
	"d7024e_group04/internal/node"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

type Handler struct {
	node node.NodeHandler
}

func NewHandler(node node.NodeHandler) *Handler {
	return &Handler{
		node: node,
	}
}

func (h *Handler) ListenAndServe(ctx context.Context) error {
	errChan := make(chan error, 1)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /objects/{hash}", h.getObject())
	mux.HandleFunc("POST /objects", h.putObject())

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
		log.Printf("http get req, %v", h.node.Me())

		hash := r.PathValue("hash")
		// check if hex value
		_, err := hex.DecodeString(hash)
		if err != nil || len(hash) != 40 { // length of sha-1
			http.Error(w, "invalid hash", http.StatusUnprocessableEntity)
			return
		}

		log.Println("request sent")

		ctx, cancel := context.WithTimeout(r.Context(), env.RPCTimeout)
		value, candidates, err := h.node.GetObject(ctx, hash)
		cancel()

		log.Println("request done")

		if err != nil {
			log.Println("err")

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if value != nil {
			log.Println("value")

			json.NewEncoder(w).Encode(value)
			return
		}

		if candidates != nil {
			log.Println("candidate")

			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(candidates)
		}

	}
}

func (h *Handler) putObject() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("http post req, %v", h.node.Me())

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		data := string(body)

		ctx, cancel := context.WithTimeout(r.Context(), env.RPCTimeout)
		hash, err := h.node.PutObject(ctx, data)
		cancel()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Location", fmt.Sprintf("/objects/%v", hash))
		w.WriteHeader(http.StatusCreated)

		response := map[string]string{
			"data": data,
		}

		json.NewEncoder(w).Encode(response)
	}
}
