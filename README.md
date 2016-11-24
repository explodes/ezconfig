*This project is mostly just an exercise, but it works really well.*

# EZConfig
Connect to databases and brokers using a well defined configuration structure.

## Example

*See /sample/sample.go for usage in practice*

```go
package main

import (
	"time"
	"log"
	"github.com/explodes/ezconfig/backoff"
	"github.com/explodes/ezconfig/opener"

	"myproject/configreader"
)

const (
        connectionRetries = 10
)

func main() {
        config := configreader.getConfig()
        connections := opener.New().
            WithRetry(connectionRetries, backoff.Exponential(10*time.Millisecond, 1*time.Second, 2)).
            WithDatabase(&config.DbConfig).
            WithProducer(&config.ProducerConfig).
            Connect()
        if connections.Err != nil {
                log.Fatalf("Error connecting, aborting: %v", connections.Err)
        }
        connections.DB.Exec(`SELECT "Hello, world!"`)
        connections.Producer.Publish("hello", "world")
}
```