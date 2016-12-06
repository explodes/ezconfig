package opener

import (
	"database/sql"
	"sync"

	"io"

	"github.com/explodes/ezconfig"
	"github.com/explodes/ezconfig/backoff"
	"github.com/explodes/ezconfig/producer"
)

// Opener is helper designed to make connecting to multiple services easy.
// It is built by specifying what to connect to (and how) and then ran.
// By default, it will only attempt to connect to its sources once each and then fail or succeed.
//
//   	connections, err := opener.New().
//   		WithRetry(connectionRetries, backoff.Constant(1*time.Second)).
//   		WithDatabase(&config.DbConfig).
//   		WithProducer(&config.ProducerConfig).
//   		Connect()
type Opener struct {
	file           string
	dbConfig       *ezconfig.DbConfig
	producerConfig *ezconfig.ProducerConfig
	retries        int
	backoff        backoff.Strategy
}

// Connections is the result of connecting to multiple sources
type Connections struct {
	DB       *sql.DB
	Producer producer.Producer
}

// New creates a New opener with no retry attempts or backoff strategy
func New() *Opener {
	return &Opener{}
}

// WithRetry sets the number of attempts to make to each source and the
// backoff strategy to utilize when re-attempting
func (co *Opener) WithRetry(retries int, strategy backoff.Strategy) *Opener {
	co.retries = retries
	co.backoff = strategy
	return co
}

// WithDatabase specifies that an attempt should be made to connect to a database
// and which settings to use to do so
func (co *Opener) WithDatabase(config *ezconfig.DbConfig) *Opener {
	co.dbConfig = config
	return co
}

// WithProducer specifies that an attempt should be made to connect to a producer
// and which settings to use to do so
func (co *Opener) WithProducer(config *ezconfig.ProducerConfig) *Opener {
	co.producerConfig = config
	return co
}

// Connect connects to the services that are set.
// In the event of error, anything successfully connected to is closed, and
// the first error received is returned.
func (co *Opener) Connect() (*Connections, error) {

	wg := sync.WaitGroup{}
	result := &Connections{}
	errs := &firstError{}

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

	if errs.err != nil {
		// close any connections we may have made
		// it is in a goroutine so that the unusuable
		// results aren't blocking the downstream
		go result.Close()
		return nil, errs.err
	}
	return result, nil
}

// connectDb connects to a database and saves the result in the given Connections
func (co *Opener) connectDb(result *Connections) error {
	database, err := InitDb(co.dbConfig, co.retries, co.backoff)
	if err != nil {
		return err
	} else {
		result.DB = database
		return nil
	}
}

// connectBroker connects to a producer and saves the result in the given Connections
func (co *Opener) connectBroker(result *Connections) error {
	prod, err := InitProducer(co.producerConfig, co.retries, co.backoff)
	if err != nil {
		return err
	} else {
		result.Producer = prod
		return nil
	}
}

// Close closes all active connections (each in independent goroutines) and returns the
// first error received
func (c *Connections) Close() error {
	return CloseAll(c.DB, c.Producer)
}

// CloseAll closes all io.Closers (each in independent goroutines) and returns the
// first error received
func CloseAll(closers ...io.Closer) error {
	wg := sync.WaitGroup{}

	errs := &firstError{}

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

// firstError is a construct that records the first error it receives via Record
type firstError struct {
	sync.Mutex
	err error
}

// Record saves the error if there has been no other error
func (e *firstError) Record(err error) {
	if err == nil {
		return
	}
	e.Lock()
	defer e.Unlock()
	if e.err == nil {
		e.err = err
	}
}
