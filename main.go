package main

import (
	"MedKick-backend/pkg/database"
	"MedKick-backend/pkg/echo"
	"MedKick-backend/pkg/echo/middleware"
	"MedKick-backend/pkg/sendgrid"
	"MedKick-backend/pkg/validator"
	"MedKick-backend/v1/cron"
	"MedKick-backend/v1/organization"
	"MedKick-backend/v1/user"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.ConnectDatabase(database.Config())
	validator.Setup()
	sendgrid.Setup()

	e := echo.Engine()

	// Add Auth
	e.Use(middleware.Auth)

	e.GET("/", echo.OnlineCheck)

	v1 := e.Group("/v1")

	// System routes
	cron.Routes(v1)

	// Main routes
	user.Routes(v1)
	organization.Routes(v1)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))))
}
