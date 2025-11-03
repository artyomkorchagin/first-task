package infrastructure

import (
	"time"

	"github.com/IBM/sarama"
	"github.com/artyomkorchagin/first-task/internal/config"
)

func ConnectConsumer(cfg config.KafkaConfig) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = cfg.ReturnErr
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second

	return sarama.NewConsumer(cfg.Brokers, config)
}
