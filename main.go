package main

import (
	"log"
	"time"

	"github.com/hoodoo-digital/sandstorm/cmd"
)

func main() {
	start := time.Now()
	cmd.Execute()
	duration := time.Since(start)
	log.Println(duration)
}
