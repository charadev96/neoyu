package server

import (
	"io/fs"
	"net/http"
	"strings"

	connectionv1 "github.com/charadev96/neoyu/gen/neoyu/connection/v1/connectionv1connect"
	"github.com/charadev96/neoyu/internal/db"
	"github.com/charadev96/neoyu/internal/handler"
	"github.com/charadev96/neoyu/internal/middleware"
	"github.com/charadev96/neoyu/internal/service"

	"connectrpc.com/connect"
	"connectrpc.com/validate"
	"github.com/rs/zerolog"
)

type DB struct {
	Connections *db.File[service.ProviderSchema]
}

type Server struct {
	db       DB
	dist     fs.FS
	provider *service.Provider
}

func New(d fs.FS, db DB) *Server {
	return &Server{
		db:       db,
		dist:     d,
		provider: service.NewProvider(db.Connections),
	}
}

func (s *Server) Serve(addr string, log *zerolog.Logger) error {
	var (
		mux = http.NewServeMux()
		api = http.NewServeMux()
	)

	path, handler := connectionv1.NewProviderServiceHandler(
		handler.NewProvider(s.provider),
		connect.WithInterceptors(
			middleware.NewLogInterceptor(log),
			validate.NewInterceptor(),
		),
	)
	api.Handle(path, handler)
	mux.Handle("/api/", http.StripPrefix("/api", api))

	dist, err := fs.Sub(s.dist, "dist")
	if err != nil {
		return err
	}
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		url := strings.TrimPrefix(r.URL.Path, "/")
		if stat, err := fs.Stat(dist, url); err == nil && !stat.IsDir() {
			http.FileServerFS(dist).ServeHTTP(w, r)
			return
		}
		http.ServeFileFS(w, r, dist, "index.html")
	})

	p := new(http.Protocols)
	p.SetHTTP1(true)
	p.SetUnencryptedHTTP2(true)
	srv := http.Server{
		Addr:      addr,
		Handler:   mux,
		Protocols: p,
	}

	return srv.ListenAndServe()
}
