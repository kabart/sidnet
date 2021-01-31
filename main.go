package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

const portDefault int = 8080

func main() {

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = portDefault
	}

	log.Printf("Starting sidnet microservice on port %d", port)
	s := newServer()

	if err := s.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Println(err)
	}
}
