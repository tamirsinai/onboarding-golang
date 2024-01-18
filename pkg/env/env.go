package env

import (
	"fmt"
	"github.com/joho/godotenv"
)

func Load() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env file")
	}
}