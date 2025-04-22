package config

import (
	"log"
	"os"

	"github.com/meilisearch/meilisearch-go"
)

var MeiliClient *meilisearch.Client

func InitMeiliSearch() {
	host := getEnv("MEILI_HOST", "localhost")
	port := getEnv("MEILI_PORT", "7700")
	apiKey := getEnv("MEILI_MASTER_KEY", "")

	MeiliClient = meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   "http://" + host + ":" + port,
		APIKey: apiKey,
	})

	_, err := MeiliClient.Health()
	if err != nil {
		log.Fatalf("❌ Cannot connect to Meilisearch at %s:%s — %v", host, port, err)
	}

	log.Println("✅ Meilisearch client is ready")
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
