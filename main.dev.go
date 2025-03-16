//go:build dev

package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/tnqbao/gau_phim_backend/config"
	"github.com/tnqbao/gau_phim_backend/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	//config.InitRedis()
	db := config.InitDB()

	ai.LoadAPIKeys()
	router := routes.SetupRouter(db)
	go vote.StartSyncJob(db)
	router.Run(":8083")
}
