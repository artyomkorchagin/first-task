package infrastructure

import (
	"github.com/IBM/sarama"
	"github.com/artyomkorchagin/first-task/internal/config"
)

func ConnectConsumer(cfg config.KafkaConfig) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = cfg.ReturnErr

	return sarama.NewConsumer(cfg.Brokers, config)
}
