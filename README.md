*This project is mostly just an exercise, but it works really well.*

# EZConfig
Connect to databases and brokers using a well defined configuration structure.

## Example

*See /sample/sample.go for a more thorough example*

```go
package main

import (
	"time"
	"log"
	"github.com/explodes/ezconfig"
	"github.com/explodes/ezconfig/backoff"
	"github.com/explodes/ezconfig/opener"
	_ "github.com/explodes/ezconfig/db/pg" // allow postgres connections
	_ "github.com/explodes/ezconfig/producer/kafka" // allow kafka connections
)

const (
        // number of attempts to make to connect to each service
        connectionRetries = 10
)

type MyConfig struct {
	ezconfig.ProducerConfig
	ezconfig.DbConfig
}

func main() {
        config := &MyConfig{}
        ezconfig.ReadConfig("local.conf", &config)
        connections, err := opener.New().
            WithRetry(connectionRetries, backoff.Exponential(10*time.Millisecond, 1*time.Second, 2)).
            WithDatabase(&config.DbConfig).
            WithProducer(&config.ProducerConfig).
            Connect()
        if err != nil {
                log.Fatalf("Error connecting, aborting: %v", err)
        }
        connections.DB.Exec(`SELECT "Hello, world!"`)
        connections.Producer.Publish("hello", "world")
}
```