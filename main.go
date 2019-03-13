package main

import (
	"microlog/settings"
	"microlog/web"
	"sync"
)

func main() {
	settings.Migrate()

	wg := &sync.WaitGroup{}

	// start web server
	wg.Add(1)
	go web.Start()

	// start inputs
	repo := settings.Inputs

	inputs, _ := repo.GetAll()

	for _, input := range inputs {
		if input.Enabled == 1 {
			listener, _ := input.GetListener()
			listener.Start()
		}
	}

	// wait all goroutines
	wg.Wait()
}
