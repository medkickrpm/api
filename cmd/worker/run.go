package worker

import (
	"MedKick-backend/pkg/database"
	"MedKick-backend/pkg/worker"
	"log"

	"github.com/joho/godotenv"
)

func Run() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.ConnectDatabase(database.Config())

	if err = worker.ProcessCPTCode99453(); err != nil {
		log.Fatalln(err.Error())
	}
}
