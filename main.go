// main.go

package main

import (
	"os"
	"github.com/joho/godotenv"
	"log"
)

// use godot package to load/read the .env file and
// return the value of the key
func getEnv(key string) string {

	// load .env file
	err := godotenv.Load(".env")
  
	if err != nil {
	  log.Fatalf("Error loading .env file")
	}
  
	return os.Getenv(key)
  }

func main() {
	a := App{}
    a.Initialize(
		getEnv("DB_USERNAME"),
		getEnv("DB_PASSWORD"),
		getEnv("DB_NAME"))

	a.Run(":8010")
}