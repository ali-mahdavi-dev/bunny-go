package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/spf13/cast"

	"bunny-go/cmd/command"
	"bunny-go/internal/user_management"
	"bunny-go/pkg/framwork/infrastructure/databases"
)

func init() {
	log.SetFlags(log.Lshortfile)
}
func main() {
	command.Execute()

	err := godotenv.Load()

	if err != nil {
		panic(err)
	}

	
	if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}

	server := gin.Default()

	user_management.Bootstrap(server, db)

	err = server.Run()
	if err != nil {
		panic(err)
	}

}

