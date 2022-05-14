package main

import (
	"log"

	"forum"
	"forum/internal/delivery"
	"forum/internal/repository"
	"forum/internal/service"
	"forum/utils"
)

func main() {
	// sqlite and html parser
	templates, err := utils.TemplateParsing("./assets/html/")
	if err != nil {
		log.Fatalf("error templates cannot be parsed: %s", err.Error())
	}

	repos := repository.NewRepository()
	services := service.NewService(repos)
	handlers := delivery.NewHandler(services, templates)
	// server
	srv := new(forum.Server)

	if err := srv.Run("8080", handlers.InitRoutes()); err != nil {
		log.Fatalf("error occured on server shutting down: %s", err.Error())
	}
}
