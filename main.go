package main

import (
	"fmt"
	"os"
	"repot/cmd/app"

	"github.com/joho/godotenv"
)

func main() {
	if _, err := os.Stat(".env"); err == nil {
		err := godotenv.Load()
		if err != nil {
			fmt.Println("Error loading .env file")
		}
	}

	app.New().Run()
}
