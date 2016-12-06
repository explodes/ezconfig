package ezconfig

import "fmt"

// Db Config is configuration in the following format:
//   [database]
//   type = "postgres"
//   host = "localhost"
//   port = 5432
//   user = "test"
//   password = "test"
//   dbname = "test"
//   ssl = "disable"
//   max_connections = 10
//
type DbConfig struct {
	Database DbHost
}

type DbHost struct {
	Type           string // only sqlite3, postgres is supported
	Host           string // file, or :memory:, for sqlite3
	Port           int
	DbName         string `toml:"dbname"`
	User           string
	Password       string
	Ssl            string
	MaxConnections int `toml:"max_connections"`
}

// ProducerConfig is config in the following format:
//   [producer]
//   type = "dummy"
//   retries = 5
//
//   [[producers]]
//   host = "docker.loc"
//   port = 9092
type ProducerConfig struct {
	Settings ProducerSettings `toml:"producer"`
	Hosts    []ProducerHost   `toml:"producers"`
}

type ProducerSettings struct {
	Type    string // "kafka" or "dummy"
	Retries int
}

type ProducerHost struct {
	Host string
	Port int
}

func (b *ProducerHost) Address() string {
	return fmt.Sprintf("%s:%d", b.Host, b.Port)
}
