package main

import "github.com/kevwan/gobench"

func main() {
	b := gobench.NewBench()
	b.Run(10000, func() {})
}
