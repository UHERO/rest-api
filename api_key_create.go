package main

import (
	"fmt"
	"github.com/UHERO/rest-api/data"
	"log"
)

func main() {
	key, err := data.CreateNewApiKey(32)
	if err != nil {
		log.Fatalf("Failure: %s\n", err)
		return
	}
	fmt.Printf("Key: |%s|\n", key)
}

