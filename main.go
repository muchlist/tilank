
package main

import (
	"github.com/joho/godotenv"
	"tilank/app"
	"log"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app.RunApp()
}