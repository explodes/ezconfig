package producer

import "fmt"

type ProducerConfig struct {
	Settings Settings `toml:"producer"`
	Hosts    []Host   `toml:"producers"`
}

type Settings struct {
	Type    string // "kafka" or "dummy"
	Retries int
}

type Host struct {
	Host string
	Port int
}

func (b *Host) Address() string {
	return fmt.Sprintf("%s:%d", b.Host, b.Port)
}

func (config *ProducerConfig) GetProducerConfig() *ProducerConfig {
	return config
}
