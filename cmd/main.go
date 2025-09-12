package main

import (
	"log"

	"github.com/ali-mahdavi-dev/bunny-go/cmd/command"
)

func init() {
	log.SetFlags(log.Lshortfile)
}
func main() {
	command.Execute()
}
