// +build ezconfig_sample

package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/explodes/ezconfig"
	"github.com/explodes/ezconfig/backoff"
	_ "github.com/explodes/ezconfig/db/pg"
	"github.com/explodes/ezconfig/opener"
	"github.com/explodes/ezconfig/producer"
	_ "github.com/explodes/ezconfig/producer/dummy"
	"github.com/explodes/jsonserv"
)

const (
	defaultConfig     = "local.conf"
	connectionRetries = 15
)

var (
	configFilePath = flag.String("config", defaultConfig, "Specify which config file to use")
)

// ServerConfig is extra configuration for our service
type ServerConfig struct {
	Host           string
	Port           int
	Debug          bool
	LogRequests    int   `toml:"log_requests"`
	MaxRequestSize int64 `toml:"max_request_size"`
}

// Config is the outermost configuration spec.
// Embedded is a ProducerConfig and DbConfig for connecting to a producer and database.
type Config struct {
	ezconfig.ProducerConfig
	ezconfig.DbConfig
	Server ServerConfig
}

// App is context we pass around to our view functions
type App struct {
	config   *Config
	db       *sql.DB
	producer producer.Producer
}

func init() {
	// set verbose logging
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime | log.Lmicroseconds)
}

// appWrap wraps our view functions as Views for jsonserv
func appWrap(f func(app *App, req *jsonserv.Request, res *jsonserv.Response)) jsonserv.View {
	return func(app interface{}, req *jsonserv.Request, res *jsonserv.Response) {
		f(app.(*App), req, res)
	}
}

// readConfig reads configuration and initializes our App's context
func readConfig() *App {
	config := &Config{}
	err := ezconfig.ReadConfig(*configFilePath, config)
	if err != nil {
		log.Fatal(err)
	}

	// connect to our producer and database with an exponential backoff strategy
	connections, err := opener.New().
		WithRetry(connectionRetries, backoff.Exponential(10*time.Millisecond, 1*time.Second, 2)).
		WithDatabase(&config.DbConfig).
		WithProducer(&config.ProducerConfig).
		Connect()

	if err != nil {
		log.Fatalf("Unable to connect: %v", err)
	}

	return &App{
		config:   config,
		db:       connections.DB,
		producer: connections.Producer,
	}
}

func main() {
	// read our configuration
	app := readConfig()
	config := app.config

	bind := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)
	log.Printf("Serving on %s", bind)

	server := jsonserv.New().
		SetApp(app).
		AddMiddleware(jsonserv.NewMaxRequestSizeMiddleware(config.Server.MaxRequestSize)).
		AddMiddleware(jsonserv.NewDebugFlagMiddleware(config.Server.Debug)).
		AddRoute(http.MethodGet, "Index", "/", appWrap(indexView)).
		AddRoute(http.MethodGet, "Error", "/error", appWrap(errorView))

	// if verbose logging is enabled, log requests as well
	if config.Server.LogRequests > 0 {
		server.AddMiddleware(jsonserv.NewLoggingMiddleware(config.Server.LogRequests > 1))
	}

	if err := server.Serve(); err != nil {
		log.Fatal(err)
	}
}

// indexView is a view demonstrating that we have a database and producer at the point of entry
func indexView(app *App, req *jsonserv.Request, res *jsonserv.Response) {
	res.Ok(map[string]interface{}{
		"hello":    "Hello, World!",
		"world":    true,
		"database": app.db,
		"producer": app.producer,
		"request":  req.String(),
	})
	app.producer.Publish("test", "hello_world")
}

// errorView is a view that simply returns a 500
func errorView(app *App, req *jsonserv.Request, res *jsonserv.Response) {
	res.Error(errors.New("failed!!!"))
}
