package main

import (
	"fmt"
	"log"

	"github.com/VincNT21/aggreGATOR/internal/config"
)

func main() {
	configData, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	err = configData.SetUser("vincnt")

	configData, err = config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	fmt.Println(configData)
}
