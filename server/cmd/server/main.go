package main

import (
	"log"

	"github.com/go-arco-admin/server/internal/config"
	"github.com/go-arco-admin/server/internal/model"
	"github.com/go-arco-admin/server/internal/router"
)

func main() {
	cfg := config.Load()
	db, err := model.Open(cfg.DBPath)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	log.Printf("listening on %s", cfg.Addr)
	if err := router.New(db, cfg).Run(cfg.Addr); err != nil {
		log.Fatal(err)
	}
}
