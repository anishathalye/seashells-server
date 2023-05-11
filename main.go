package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/anishathalye/seashells-server/config"
	"github.com/anishathalye/seashells-server/datamanager"
	"github.com/caarlos0/env"
)

const gcTime = 1 * 24 * time.Hour
const dataLimit = 250 * 1024 // data limit in bytes
const perUserLimit = 5
const idLength = 8
const writeTimeout = 10 * time.Second
const pingTimeout = 30 * time.Second
const pingPeriod = (pingTimeout * 9) / 10
const maxMessageSize = 100

func main() {

	cfg := config.Config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("failed to parse env vars: %v\n", err)
		return
	}

	fmt.Printf("%+v\n", cfg)

	rand.Seed(time.Now().UTC().UnixNano())

	manager := datamanager.New(dataLimit, perUserLimit, gcTime)

	go runNetcatServer(manager)

	runWeb(manager)

}
