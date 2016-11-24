package opener

import (
	"database/sql"
	"sync"

	"github.com/explodes/ezconfig/backoff"
	"github.com/explodes/ezconfig/db"
	"github.com/explodes/ezconfig/producer"
)

type Opener struct {
	file           string
	dbConfig       *db.DbConfig
	producerConfig *producer.ProducerConfig
	retries        int
	backoff        backoff.Strategy
}

type ConnectionResult struct {
	DB       *sql.DB
	Producer producer.Producer
	Err      error
}

func New() *Opener {
	return &Opener{}
}

func (co *Opener) WithRetry(retries int, strategy backoff.Strategy) *Opener {
	co.retries = retries
	co.backoff = strategy
	return co
}

func (co *Opener) WithDatabase(config *db.DbConfig) *Opener {
	co.dbConfig = config
	return co
}

func (co *Opener) WithProducer(config *producer.ProducerConfig) *Opener {
	co.producerConfig = config
	return co
}

func (co *Opener) Connect() *ConnectionResult {

	wg := sync.WaitGroup{}
	result := &ConnectionResult{}

	if co.dbConfig != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			co.connectDb(result)
		}()
	}

	if co.producerConfig != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			co.connectBroker(result)
		}()
	}

	wg.Wait()

	return result
}

func (co *Opener) connectDb(result *ConnectionResult) {
	database, err := db.InitDb(co.dbConfig, co.retries, co.backoff)
	if err != nil {
		result.Err = err
	} else {
		result.DB = database
	}
}

func (co *Opener) connectBroker(result *ConnectionResult) {
	prod, err := producer.InitProducer(co.producerConfig, co.retries, co.backoff)
	if err != nil {
		result.Err = err
	} else {
		result.Producer = prod
	}
}
