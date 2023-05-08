package main

import (
	"math/rand"
	"time"

	"github.com/kevwan/gobench"
)

func main() {
	b := gobench.NewBenchWithConfig(gobench.Config{
		Title: "sleep",
		Host:  "localhost",
		Port:  8282,
	})
	b.Run(10000, func() {
		n := rand.Intn(100)
		time.Sleep(time.Millisecond * time.Duration(n))
	})
}
