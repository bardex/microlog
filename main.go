package main

import (
	"microlog/web"
	"sync"
	"microlog/settings"
)

func main() {
	wg := &sync.WaitGroup{}

	// start web server
	wg.Add(1)
	go web.Start()

	// start inputs
	repo := settings.Inputs{}
	inputs, _ := repo.GetAll()



	// wait all goroutines
	wg.Wait()
}
