package main

import (
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load("conf.example.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	BootstrapAndRun()
}
