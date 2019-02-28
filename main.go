package main

import (
	"microlog/settings"
	"microlog/web"
	"sync"
)

func main() {
	wg := &sync.WaitGroup{}

	// start web server
	wg.Add(1)
	go web.Start()

	// start inputs
	repo := settings.Inputs

	//repo.Install()

	inputs, _ := repo.GetAll()

	for _, input := range inputs {
		go input.GetListener().Start()
	}

	// wait all goroutines
	wg.Wait()
}
