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

	eventRepo := repository.NewEvenRepositoryImpl(db)
	eventService := service.NewEventServiceImpl(eventRepo, validate)
	eventHandler := handler.NewEventHandlerImpl(eventService)

	ticketRepo := repository.NewTicketRepositoryImpl(db)
	ticketService := service.NewTicketServiceImpl(ticketRepo, eventRepo, validate)
	ticketHandler := handler.NewTicketHandlerImpl(ticketService)

	routes := routes.NewRouter(userHandler, eventHandler, ticketHandler)

	port := os.Getenv("PORT_APP")
	routes.Run(port)
	log.Print("server run in port ", port)
}