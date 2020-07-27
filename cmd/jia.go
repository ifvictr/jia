package main

import (
	"fmt"

	"github.com/ifvictr/jia/pkg/jia"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func main() {
	fmt.Println("Starting Jia…")
	config := jia.NewConfig()

	// Start receiving messages
	fmt.Println(fmt.Sprintf("Listening on port %d", config.Port))
	jia.StartServer(config)
}
