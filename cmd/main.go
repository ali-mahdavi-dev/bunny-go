package main

import (
	"log"

	"shikposh-backend/cmd/commands"
)

func main() {
	commands.Execute()
}

func init() {
	log.SetFlags(log.Lshortfile)
}
