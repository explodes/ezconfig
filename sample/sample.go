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

type ServerConfig struct {
	Host           string
	Port           int
	Debug          bool
	LogRequests    int   `toml:"log_requests"`
	MaxRequestSize int64 `toml:"max_request_size"`
}

type Config struct {
	ezconfig.ProducerConfig
	ezconfig.DbConfig
	Server ServerConfig
}

type App struct {
	config   *Config
	db       *sql.DB
	producer producer.Producer
}

func init() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime | log.Lmicroseconds)
}

func appWrap(f func(app *App, req *jsonserv.Request, res *jsonserv.Response)) jsonserv.View {
	return func(app interface{}, req *jsonserv.Request, res *jsonserv.Response) {
		f(app.(*App), req, res)
	}
}

func readConfig() *App {
	config := &Config{}
	err := ezconfig.ReadConfig(*configFilePath, config)
	if err != nil {
		log.Fatal(err)
	}

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
	app := readConfig()
	config := app.config

	bind := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)
	log.Printf("Serving on %s", bind)

	server := jsonserv.New()
	if config.Server.LogRequests > 0 {
		server.AddMiddleware(jsonserv.NewLoggingMiddleware(config.Server.LogRequests > 1))
	}
	err := server.SetApp(app).
		AddMiddleware(jsonserv.NewMaxRequestSizeMiddleware(config.Server.MaxRequestSize)).
		AddMiddleware(jsonserv.NewDebugFlagMiddleware(config.Server.Debug)).
		AddRoute(http.MethodGet, "Index", "/", appWrap(indexView)).
		AddRoute(http.MethodGet, "Error", "/error", appWrap(errorView)).
		Serve(bind)
	if err != nil {
		log.Fatal(err)
	}
}

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

func errorView(app *App, req *jsonserv.Request, res *jsonserv.Response) {
	res.Error(errors.New("failed!!!"))
}
