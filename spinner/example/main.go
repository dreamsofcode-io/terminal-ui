package main

import (
	"log"
	"time"

	"github.com/dreamsofcode-io/terminal-ui/spinner"
)

func main() {
	s := spinner.New(spinner.Config{})

	log.Println("Starting the spinner")
	s.Start()

	time.Sleep(time.Second * 5)
	s.Stop()

	log.Println("Spinner stopped")
}
