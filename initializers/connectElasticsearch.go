package initializers

import (
	"github.com/elastic/go-elasticsearch/v8"
	"log"
	"os"
)

var ES *elasticsearch.TypedClient

func ConnectElasticsearch() {
	var err error
	url := os.Getenv("ES_HOST")
	username := os.Getenv("ES_USERNAME")
	password := os.Getenv("ES_PASSWORD")
	cfg := elasticsearch.Config{
		Addresses: []string{url},
		Username:  username,
		Password:  password,
	}
	ES, err = elasticsearch.NewTypedClient(cfg)
	if err != nil {
		log.Fatalf("Error getting when creating new ES client: %s", err)
	}
}

func CheckElasticIndex() {
	_, err := ES.Indices.Create("item").Do(nil)
	if err != nil {
		log.Printf("Error getting when ES createing new index: %s", err)
	}
}
