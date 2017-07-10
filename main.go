package main

import (
	"math/rand"
	"github.com/anishathalye/seashells-server/datamanager"
	"time"
)

const baseUrl = "https://seashells.io/v/"
const gcTime = 1 * 24 * time.Hour
const dataLimit = 250 * 1024 // data limit in bytes
const perUserLimit = 5
const idLength = 8
const writeTimeout = 10 * time.Second
const pingTimeout = 30 * time.Second
const pingPeriod = (pingTimeout * 9) / 10
const maxMessageSize = 100

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	manager := datamanager.New(dataLimit, perUserLimit, gcTime)
	go runNetcatServer(manager)
	runWeb(manager)
}
