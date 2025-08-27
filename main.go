package main

import (
	"log"
	"os"
	"ticket/config"
	"ticket/handler"
	"ticket/repository"
	"ticket/routes"
	"ticket/service"

	"github.com/go-playground/validator/v10"
)

func main() {
	db := config.Db()
	validate := validator.New()

	userRepo := repository.NewUserReposiitoryImpl(db)
	userService := service.NewUserServiceImpl(userRepo, validate)
	userHandler := handler.NewUserHandlerImpl(userService)


	routes := routes.NewRouter(userHandler)

	port := os.Getenv("PORT_APP")
	routes.Run(port)
	log.Print("server run in port ", port)
}