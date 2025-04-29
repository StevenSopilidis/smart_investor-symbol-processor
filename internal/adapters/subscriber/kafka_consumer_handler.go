package subscriber

import (
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
	"github.com/stevensopi/smart_investor/symbol_processor/internal/core/domain"
	"github.com/stevensopi/smart_investor/symbol_processor/internal/core/ports"
)

type KafkaConsumerHandler struct {
	repo ports.ISymbolRepo
}

func NewKafkaConsumerHandler(repo ports.ISymbolRepo) *KafkaConsumerHandler {
	return &KafkaConsumerHandler{
		repo: repo,
	}
}

func (h *KafkaConsumerHandler) Setup(session sarama.ConsumerGroupSession) error {
	log.Println("[KafkaConsumerHandler] Consumer setup: initializing...")
	return nil
}

func (h *KafkaConsumerHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	log.Println("[KafkaConsumerHandler] Consumer cleanup: releasing resources...")
	return nil
}

func (h *KafkaConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		data := msg.Value

		var symbol_data domain.SymbolData

		err := json.Unmarshal(data, &symbol_data)
		if err != nil {
			log.Println("---> Could not unmarshal data")
			continue
		}

		log.Printf("---> Received data for ticker %s\n", symbol_data.Ticker)

		h.repo.Put(symbol_data)
	}
	return nil
}
