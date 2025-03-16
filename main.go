//go:build !dev

package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/tnqbao/gau_phim_backend/config"
	"github.com/tnqbao/gau_phim_backend/routes"
)

func main() {
	err := godotenv.Load("/gau_phim/.env.")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	config.InitRedis()

	db := config.InitDB()

	router := routes.SetupRouter(db)

	router.Run(":8083")
}
