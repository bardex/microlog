package main

import (
	"microlog/input"
	"microlog/web"
	"sync"
)

func main() {

	wg := &sync.WaitGroup{}

	// start web server
	wg.Add(1)
	go web.Start()

	// start inputs
	i, _ := input.Create("udp", ":8081")
	go i.Start()

	// wait all goroutines
	wg.Wait()
}
