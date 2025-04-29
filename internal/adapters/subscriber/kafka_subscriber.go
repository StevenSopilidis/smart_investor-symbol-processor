package subscriber

import (
	"context"
	"log"

	"github.com/IBM/sarama"
	"github.com/stevensopi/smart_investor/symbol_processor/internal/adapters/config"
	"github.com/stevensopi/smart_investor/symbol_processor/internal/core/ports"
)

type KafkaSubscriber struct {
	consumerGroup sarama.ConsumerGroup
	handler       *KafkaConsumerHandler
	topic         string
}

func NewKafkaSubscriber(config config.Config, repo ports.ISymbolRepo) (ports.ISubscriber, error) {
	sconfig := sarama.NewConfig()
	sconfig.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	sconfig.Consumer.Offsets.Initial = sarama.OffsetNewest

	brokers := []string{config.KafkaBroker}
	consumerGroup, err := sarama.NewConsumerGroup(brokers, config.KafkaGroupId, sconfig)
	if err != nil {
		return nil, err
	}

	handler := NewKafkaConsumerHandler(repo)

	return &KafkaSubscriber{
		consumerGroup: consumerGroup,
		handler:       handler,
		topic:         config.SymbolTopic,
	}, nil
}

func (s *KafkaSubscriber) Run() {
	for {
		if err := s.consumerGroup.Consume(context.Background(), []string{s.topic}, s.handler); err != nil {
			log.Printf("Error consuming messages: %v\n", err)
		}
	}
}
