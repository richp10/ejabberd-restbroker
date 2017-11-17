package main

import (
	"runtime"
	"strconv"

	"ejabberd-restbroker/lib/boot"

	"ejabberd-restbroker/lib/flight"

	"ejabberd-restbroker/job/process"

	"github.com/jrallison/go-workers"
)


func main() {
	boot.RegisterServices()
	c := flight.Context().Config

	// Configure the worker pool
	workers.Configure(map[string]string{
		"server":    c.RedisHost,
		"database":  strconv.Itoa(c.RedisDatabase),
		"pool":      strconv.Itoa(c.RedisPoolCount),
		"process":   c.ProcessID,
		"namespace": c.RedisNamespace,
	})

	// pull messages from named queue with concurrency of 10
	// These are then passed to the process.Job within a
	// goroutine
	workers.Process(c.RestQueue, process.Job, c.Concurreny)

	// By default stats will be available at http://localhost:8080/stats
	if c.Stats {
		go workers.StatsServer(c.StatsPort)
	}

	// Blocks until process is told to exit via unix signal
	workers.Run()
}

func init() {
	// Use all CPU cores
	runtime.GOMAXPROCS(runtime.NumCPU())
}
