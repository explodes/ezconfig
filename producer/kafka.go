package producer

import (
	"errors"
	"strconv"
	"time"

	"github.com/Shopify/sarama"
)

func kafkaValidateConfig(conf *ProducerConfig) error {
	if conf.Hosts == nil || len(conf.Hosts) == 0 {
		return errors.New("Invalid producer configration: No [[producers]] entry in configuration")
	}
	return nil
}

func kafkaInitProducer(conf *ProducerConfig) (Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Retry.Max = conf.Settings.Retries
	config.Producer.RequiredAcks = sarama.WaitForAll
	producers := []string{}
	for _, producer := range conf.Hosts {
		producers = append(producers, producer.Address())
	}
	producer, err := sarama.NewAsyncProducer(producers, config)
	if err != nil {
		return nil, err
	}
	kafka := kafkaProducer{
		p:    producer,
		conf: conf,
	}
	return &kafka, nil
}

type kafkaProducer struct {
	p    sarama.AsyncProducer
	conf *ProducerConfig
}

func (k kafkaProducer) Publish(topic string, message string) {
	strTime := strconv.Itoa(int(time.Now().Unix()))
	k.p.Input() <- &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(strTime),
		Value: sarama.StringEncoder(message),
	}
}

func (k kafkaProducer) Close() error {
	k.p.Close()
	return nil
}
