package server

import (
	"MedKick-backend/pkg/database"
	"MedKick-backend/pkg/echo"
	"MedKick-backend/pkg/echo/middleware"
	"MedKick-backend/pkg/s3"
	"MedKick-backend/pkg/sendgrid"
	"MedKick-backend/pkg/validator"
	"MedKick-backend/v1/careplan"
	"MedKick-backend/v1/cron"
	"MedKick-backend/v1/device"
	"MedKick-backend/v1/interaction"
	"MedKick-backend/v1/organization"
	"MedKick-backend/v1/user"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
)

func Run() {
	err := godotenv.Load(".env")
	if err != nil && os.Getenv("ENV") != "production" {
		log.Fatalf("Error loading .env file.")
	}

	database.ConnectDatabase(database.Config())
	validator.Setup()
	sendgrid.Setup()
	s3.Setup()

	middleware.Setup()

	e := echo.Engine()

	// Add Swagger
	echo.Swagger(e)

	// Add Auth
	// TODO - Performance task for Raaj later
	e.Use(middleware.Auth)

	e.GET("/", echo.OnlineCheck)

	v1 := e.Group("/v1")
	v1.GET("/", echo.OnlineCheck)

	// System routes
	cron.Routes(v1)

	// Main routes
	user.Routes(v1)
	organization.Routes(v1)
	interaction.Routes(v1)
	careplan.Routes(v1)
	device.Routes(v1)

	go func() {
		if err := e.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))); err != nil && !errors.Is(err, http.ErrServerClosed) {
			e.Logger.Fatal("Shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
