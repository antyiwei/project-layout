package main

import (
	"log"

	"project-layout/internal/router"
	"project-layout/internal/user/handler"
	"project-layout/internal/user/repository"
	"project-layout/internal/user/service"
)

func main() {
	userRepo := &repository.Repo{}
	userSvc := &service.Service{Repo: userRepo}
	userH := &handler.Handler{Service: userSvc}

	r := router.Engine(router.Handlers{User: userH})
	log.Println("listening on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
