package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/vvv/api/http"
	"gopkg.in/yaml.v2"
)

func main() {
	flag.Parse()

	confFile := flag.Arg(0)
	addr := flag.Arg(1)

	conf, err := LoadConfig(confFile)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	HTTPServer := http.NewServer(conf.LB.Addr)

	// Copy configuration settings to the HTTP server.
	HTTPServer.Addr = addr

	// Start the HTTP server.
	if err := HTTPServer.Open(); err != nil {
		panic(err)
	}

	fmt.Printf("api server start at %s\n", addr)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	fmt.Printf("api server start to shutdown on signal, %v\n", <-c)
	cancel()

	<-ctx.Done()

	fmt.Println("api server shutdown")
}

// Config represents the CLI configuration file.
type Config struct {
	LB struct {
		Addr string `yaml:"addr"`
	} `yaml:"lb"`
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
