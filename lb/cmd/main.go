package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/vvv/lb/http"
	"github.com/vvv/lb/redis"
	"github.com/vvv/lb/wrr"
	"gopkg.in/yaml.v2"
)

func main() {
	flag.Parse()

	conf, err := LoadConfig(flag.Arg(0))
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	HTTPServer := http.NewServer()

	// Copy configuration settings to the HTTP server.
	HTTPServer.Addr = conf.HTTP.Addr

	// Attach underlying services to the HTTP server.
	redisCli, err := redis.NewClient(ctx, conf.Redis)
	if err != nil {
		panic(err)
	}
	selectService := redis.NewSelectService(redisCli, wrr.NewMinWRR())
	selectService.BackgroundUpdate(ctx)

	podService := redis.NewPodService(redisCli)
	podService.BackgroundUpdate(ctx)

	HTTPServer.SelectService = selectService
	HTTPServer.PodService = podService

	// Start the HTTP server.
	if err := HTTPServer.Open(); err != nil {
		panic(err)
	}

	fmt.Println("lb server start")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	fmt.Printf("lb server start to shutdown on signal, %v\n", <-c)
	cancel()

	<-ctx.Done()

	fmt.Println("lb server shutdown")
}

// Config represents the CLI configuration file.
type Config struct {
	HTTP struct {
		Addr string `yaml:"addr"`
	} `yaml:"http"`
	Redis redis.Config `yaml:"redis"`
}

func LoadConfig(filename string) (*Config, error) {
	raw, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = yaml.Unmarshal(raw, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
