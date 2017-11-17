// Package flight provides access to the application settings safely.
// Derived from MIT Licensed code from Copyright (c) 2016 Blue Jay
package flight

import (
	"net/http"
	"sync"

	"ejabberd-restbroker/lib/env"
)

var (
	httpClient *http.Client
	configInfo env.Info
	mutex      sync.RWMutex
)

// Info structure for the application settings.
type Info struct {
	HttpClient *http.Client
	Config     env.Info
}

// Store http client - enabling tcp connection reuse in goroutines
func StoreHttpClient(ht *http.Client) {
	mutex.Lock()
	httpClient = ht
	mutex.Unlock()
}

// Store Config - so controllers can access
func StoreConfig(ci env.Info) {
	mutex.Lock()
	configInfo = ci
	mutex.Unlock()
}

// Context returns the application settings.
func Context() Info {

	i := Info{
		HttpClient: httpClient,
		Config:     configInfo,
	}

	return i
}
