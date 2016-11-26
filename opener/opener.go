package opener

import (
	"database/sql"
	"sync"

	"github.com/explodes/ezconfig/backoff"
	"github.com/explodes/ezconfig/db"
	"github.com/explodes/ezconfig/producer"
	"io"
)

type Opener struct {
	file           string
	dbConfig       *db.DbConfig
	producerConfig *producer.ProducerConfig
	retries        int
	backoff        backoff.Strategy
}

type Connections struct {
	DB       *sql.DB
	Producer producer.Producer
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

func (co *Opener) Connect() (*Connections, error) {

	wg := sync.WaitGroup{}
	result := &Connections{}
	errs := &lastError{}

	if co.dbConfig != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			errs.Record(co.connectDb(result))
		}()
	}

	if co.producerConfig != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			errs.Record(co.connectBroker(result))
		}()
	}

	wg.Wait()

	return result, errs.err
}

func (co *Opener) connectDb(result *Connections) error {
	database, err := db.InitDb(co.dbConfig, co.retries, co.backoff)
	if err != nil {
		return err
	} else {
		result.DB = database
		return nil
	}
}

func (co *Opener) connectBroker(result *Connections) error {
	prod, err := producer.InitProducer(co.producerConfig, co.retries, co.backoff)
	if err != nil {
		return err
	} else {
		result.Producer = prod
		return nil
	}
}

func (c *Connections) Close() error {
	return CloseAll(c.DB, c.Producer)
}

func CloseAll(closers ...io.Closer) error {
	wg := sync.WaitGroup{}

	errs := &lastError{}

	closeAndRecordError := func(c io.Closer) {
		if c == nil {
			return
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			errs.Record(c.Close())
		}()
	}

	for _, closer := range closers {
		closeAndRecordError(closer)
	}

	wg.Wait()
	return errs.err
}

type lastError struct {
	sync.Mutex
	err error
}

func (e *lastError) Record(err error) {
	if err == nil {
		return
	}
	e.Lock()
	defer e.Unlock()
	if e.err == nil {
		e.err = err
	}
}
