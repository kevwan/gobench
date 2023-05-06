package main

import (
	"math/rand"
	"time"

	"github.com/kevwan/gobench"
)

func main() {
	b := gobench.NewBench()
	b.Run(10000, func() {
		n := rand.Intn(100)
		time.Sleep(time.Millisecond * time.Duration(n))
	})
}
