package main

import (
	"log"

	"github.com/ali-mahdavi-dev/bunny-go/cmd/command"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

// @title						bunny-go API Documentation.
// @version					1.0.0
// @description				API documentation for Bunny-go levels
//
// @description				توضیح: فلو اندپوینت های احراز هویت سمت کاربر
// @schemes					http https
// @securityDefinitions.apikey	BearerAuth
// @type						apiKey
// @in							header
// @name						Authorization
// @description				Type "Bearer" followed by a space and JWT token.
func main() {
	command.Execute()
}
