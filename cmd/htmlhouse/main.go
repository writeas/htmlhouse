package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/writeas/htmlhouse"
	"os"
)

func main() {
	source := os.Getenv("CONFIG_SOURCE")
	if source == "" {
		source = ".env"
	}
	err := godotenv.Load(source)
	if err != nil {
		fmt.Println("unable to load configuration file", source, ":", err)
	}

	htmlhouse.Serve()
}
