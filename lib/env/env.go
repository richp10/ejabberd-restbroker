// Package env reads the application settings.
// Derived from MIT Licensed code from Copyright (c) 2016 Blue Jay
package env

import (
	"encoding/json"
	"net"
	"net/url"

	"ejabberd-restbroker/lib/jsonconfig"
)

// Info structure for the application settings.
type Info struct {
	MaxIdleConnections int
	RequestTimeout     int
	RedisHost          string
	RedisDatabase      int
	RedisPoolCount     int
	RedisNamespace     string
	ProcessID          string
	Token              string
	JabberDomain       string
	RestUser           string
	RestPass           string
	RestURL            string
	RestQueue          string
	ResponseQueue      string
	Concurreny         int
	Stats              bool
	StatsPort          int
	JabberRestHost     string
	path               string
}

// Path returns the env.json path
func (c *Info) Path() string {
	return c.path
}

// ParseJSON unmarshals bytes to structs
func (c *Info) ParseJSON(b []byte) error {
	return json.Unmarshal(b, &c)
}

// New returns a instance of the application settings.
func New(path string) *Info {
	return &Info{
		path: path,
	}
}

// LoadConfig reads the configuration file.
func LoadConfig(configFile string) (*Info, error) {
	// Create a new configuration with the path to the file
	config := New(configFile)

	// Load the configuration file
	err := jsonconfig.Load(configFile, config)
	if err != nil {
		panic(err)
	}

	// Check and set default values..
	if config.MaxIdleConnections == 0 {
		config.MaxIdleConnections = 20
	}
	if config.RequestTimeout == 0 {
		config.RequestTimeout = 10
	}

	if config.RedisHost == "" {
		config.RedisHost = "localhost:6379"
	}

	if config.RedisPoolCount == 0 {
		config.RedisPoolCount = 15
	}

	if config.ProcessID == "" {
		config.ProcessID = "1"
	}

	if config.RestURL == "" {
		config.RestURL = "localhost:5281"
	}

	if config.Concurreny == 0 {
		config.Concurreny = 20
	}

	if config.StatsPort == 0 {
		config.StatsPort = 8080
	}

	// Extract the host of the rest - might not be the same as the domain..
	u, err := url.Parse(config.RestURL)
	if err != nil {
		panic(err)
	}
	config.JabberRestHost, _, _ = net.SplitHostPort(u.Host)

	// Return the configuration
	return config, err
}
