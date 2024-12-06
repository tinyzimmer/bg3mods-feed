package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/tinyzimmer/bg3mods-feed/internal/config"
	"github.com/tinyzimmer/bg3mods-feed/internal/feed"
)

type Server struct {
	srv *http.Server
}

type ServerOptions struct {
	Generator feed.Generator
	Addr      string
}

func NewServer(opts ServerOptions) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /feed", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		data, err := opts.Generator.GetFeed(r.Context(), feed.OptionsFromQuery(r.URL))
		if err != nil {
			http.Error(w, "Failed to generate feed: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if data.Format == config.FormatJSON {
			w.Header().Set("Content-Type", "application/json")
		} else {
			w.Header().Set("Content-Type", "application/xml")
		}
		w.Header().Set("X-Feed-Generation-Time", time.Since(start).String())
		fmt.Fprint(w, string(data.Content))
	})
	return &Server{
		srv: &http.Server{
			Addr:    opts.Addr,
			Handler: logRequests(mux),
		},
	}
}

func (s *Server) ListenAndServe() error {
	log.Println("Listening on", s.srv.Addr)
	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func logRequests(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h.ServeHTTP(w, r)
		log.Println(r.Method, r.URL.String(), "-", time.Since(start).String())
	})
}
