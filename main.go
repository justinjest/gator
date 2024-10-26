package main

import (
	"fmt"
	"log"

	json_parser "github.com/justinjest/gator/internal/config"
)

func main() {
	config, err := json_parser.Read()
	if err != nil {
		log.Fatalf("Error %v", err)
	}
	_, err = json_parser.SetUser(config, "Jessica")
	if err != nil {
		log.Fatalf("Error %v", err)
	}
	newConfig, err := json_parser.Read()
	if err != nil {
		log.Fatalf("Error %v", err)
	}
	fmt.Printf("%v\n", newConfig)
}
