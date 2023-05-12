package main

import (
	"errors"
	"math/rand"
	"time"

	"github.com/kevwan/gobench"
)

func main() {
	b := gobench.NewBenchWithConfig(gobench.Config{
		Title:    "sleep",
		Host:     "localhost",
		Port:     8282,
		Duration: time.Minute,
	})
	b.RunErr(1000, func() error {
		n := rand.Intn(10)
		time.Sleep(time.Millisecond * time.Duration(n))
		if n < 3 {
			return errors.New("error")
		}
		return nil
	})
}
