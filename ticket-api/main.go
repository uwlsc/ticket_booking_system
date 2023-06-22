package main

import (
	"ticket-api/cmd"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	cmd.Execute()
}
