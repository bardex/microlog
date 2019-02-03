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
	input.StartAll()

	// wait all goroutines
	wg.Wait()
}
