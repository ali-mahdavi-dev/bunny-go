package main

import (
	"log"

	"bunny-go/cmd/command"
)

func init() {
	log.SetFlags(log.Lshortfile)
}
func main() {
	command.Execute()
}
