package main

import (
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("conf.example.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	BootstrapAndRun()
}
