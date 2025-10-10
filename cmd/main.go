package main

import (
	"log"

	"github.com/ali-mahdavi-dev/bunny-go/cmd/commands"
)

func main() {
	commands.Execute()
}

func init() {
	log.SetFlags(log.Lshortfile)
}
