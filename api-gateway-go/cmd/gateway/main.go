package main

import (
	"log"
	"net/http"

	"api-gateway-go/internal/config"
	"api-gateway-go/internal/router"
	"api-gateway-go/pkg/logger"
)

func main() {
	cfg := config.Load("configs/routes.yaml")

	logg := logger.New()

	r := router.NewRouter(cfg, logg)

	log.Println("Gateway running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
