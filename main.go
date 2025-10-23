package main

import (
	"projgen/cmd"

	_ "projgen/cmd/react"
	_ "projgen/cmd/spring"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	cmd.Execute()
}
