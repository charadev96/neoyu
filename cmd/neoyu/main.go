package main

import (
	"embed"

	connectionv1 "github.com/charadev96/neoyu/gen/neoyu/connection/v1"
	"github.com/charadev96/neoyu/internal/db"
	"github.com/charadev96/neoyu/internal/middleware"
	"github.com/charadev96/neoyu/internal/server"
)

//go:embed all:dist/*
var dist embed.FS
var addr = "127.0.0.1:8080"

func main() {
	log := middleware.NewConsoleLogger()

	svr := server.New(dist, server.DB{
		Connection: db.NewFile[*connectionv1.ConnectionStore]("data/connections.yaml"),
	})

	log.Info().
		Str("on", addr).
		Msg("serving")

	if err := svr.Serve(addr, &log); err != nil {
		log.Fatal().
			Err(err).
			Msg("failed")
	}
}
