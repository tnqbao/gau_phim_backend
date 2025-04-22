//go:build prod

package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/tnqbao/gau_phim_backend/config"
	"github.com/tnqbao/gau_phim_backend/routes"
)

func main() {
	err := godotenv.Load("/gau_phim/.env.flix")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	config.InitRedis()
	config.InitMeiliSearch()
	db := config.InitDB()

	router := routes.SetupRouter(db)

	router.Run(":8083")
}
