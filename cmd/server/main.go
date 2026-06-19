package main

import (
	"log"

	"project-layout/internal/router"
)

func main() {
	// TODO: load config, open db, wire domain services, then register routes.
	r := router.Engine(router.Handlers{})
	log.Println("listening on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
