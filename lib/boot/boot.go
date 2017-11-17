package boot

import (
	"ejabberd-restbroker/lib/env"
	"ejabberd-restbroker/lib/flight"
	"log"
	"net/http"
	"time"
)

// RegisterServices sets up all the web components.
func RegisterServices() {

	// Load the configuration file and store in flight
	config, err := env.LoadConfig("env.json")
	if err != nil {
		log.Fatalln(err)
	}
	flight.StoreConfig(*config)

	// Create and store the http connection so it can be reused.
	httpClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: config.MaxIdleConnections,
		},
		Timeout: time.Duration(config.RequestTimeout) * time.Second,
	}
	flight.StoreHttpClient(httpClient)

}
