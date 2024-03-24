package main

import (
	"context"
	"log"
	"net/http"

	"github.com/vandenbill/social-media-10k-rps/internal/cfg"
	"github.com/vandenbill/social-media-10k-rps/internal/handler"
	"github.com/vandenbill/social-media-10k-rps/internal/repo"
	"github.com/vandenbill/social-media-10k-rps/internal/service"
	"github.com/vandenbill/social-media-10k-rps/pkg/env"
	"github.com/vandenbill/social-media-10k-rps/pkg/postgre"
	"github.com/vandenbill/social-media-10k-rps/pkg/router"
	"github.com/vandenbill/social-media-10k-rps/pkg/validator"
)

func main() {
	env.LoadEnv()

	ctx := context.Background()
	router := router.NewRouter()
	conn := postgre.GetConn(ctx)
	defer conn.Close()
	validator := validator.New()

	cfg := cfg.Load()
	repo := repo.NewRepo(conn)
	service := service.NewService(repo, validator, cfg)
	handler.NewHandler(router, service, cfg)

	log.Println("server started on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalln("fail start server:", err)
	}
}
