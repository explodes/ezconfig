package kafka

import (
	"errors"
	"strconv"
	"time"

	"github.com/Shopify/sarama"
	"github.com/explodes/ezconfig"
	"github.com/explodes/ezconfig/producer"
	"github.com/explodes/ezconfig/producer/registry"
)

const (
	// kafkaProducerType is the value to use in configuration to connect to this producer type
	kafkaProducerType = "kafka"
)

// init registers the init and validation functions with the registry
func init() {
	registry.Register(kafkaProducerType, initProducer, validateConfig)
}

// validateConfig makes sure all the required settings are present for the database
func validateConfig(conf *ezconfig.ProducerConfig) error {
	if conf.Hosts == nil || len(conf.Hosts) == 0 {
		return errors.New("Invalid producer configration: No [[producers]] entry in configuration")
	}
	return nil
}

// initProducer establishes a connection with the given configuration
func initProducer(conf *ezconfig.ProducerConfig) (producer.Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Retry.Max = conf.Settings.Retries
	config.Producer.RequiredAcks = sarama.WaitForAll
	producers := []string{}
	for _, p := range conf.Hosts {
		producers = append(producers, p.Address())
	}
	p, err := sarama.NewAsyncProducer(producers, config)
	if err != nil {
		return nil, err
	}
	kafka := kafkaProducer{
		p:    p,
		conf: conf,
	}
	return &kafka, nil
}

// kafkaProducer publishes messages to kafka
type kafkaProducer struct {
	p    sarama.AsyncProducer
	conf *ezconfig.ProducerConfig
}

// Publish publishes messages to kafka
func (k kafkaProducer) Publish(topic string, message string) {
	strTime := strconv.Itoa(int(time.Now().Unix()))
	k.p.Input() <- &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(strTime),
		Value: sarama.StringEncoder(message),
	}
}

// Close closes the connection to kafka
func (k kafkaProducer) Close() error {
	k.p.Close()
	return nil
}
