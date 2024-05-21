package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("conf.example.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	fmt.Println("ACCESS_TOKEN_DURATION:", os.Getenv("ACCESS_TOKEN_DURATION"))

	BootstrapAndRun()
}
