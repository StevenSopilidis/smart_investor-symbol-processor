package main

import (
	"log"
	"sync"

	"github.com/stevensopi/smart_investor/symbol_processor/internal/adapters/config"
	"github.com/stevensopi/smart_investor/symbol_processor/internal/adapters/repo"
	"github.com/stevensopi/smart_investor/symbol_processor/internal/adapters/subscriber"
)

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("---> Could not load config %v\n", err)
	}

	repo, err := repo.NewRedisSymbolRepo(config)
	if err != nil {
		log.Fatalf("---> Could not create repo %v\n", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		subscriber, err := subscriber.NewKafkaSubscriber(config, repo)
		if err != nil {
			log.Fatalf("---> Could not create subscriber %v\n", err)
		}
		subscriber.Run()
	}()

}
