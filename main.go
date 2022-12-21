package main

import (
	"log"
	"time"

	"github.com/connorwade/dachshund/cmd"
)

func main() {
	start := time.Now()
	cmd.Execute()
	duration := time.Since(start)
	log.Println(duration)
}
