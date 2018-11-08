package main

import (
	"fmt"
	"microlog/web"
	"time"
)

func main() {
	fmt.Println("test")

	go web.Start()

	fmt.Println("ddd")

	time.Sleep(60 * time.Minute)
}