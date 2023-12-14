package gsheet

import (
	"context"
	"encoding/base64"
	"log"
	"os"

	"golang.org/x/oauth2/google"
	"gopkg.in/Iwark/spreadsheet.v2"
)

var Service *spreadsheet.Service

func Setup() {
	GSHEET_SECRET_BASE64 := os.Getenv("GSHEET_SECRET")

	// Decode from base64
	GSHEET_SECRET, err := base64.StdEncoding.DecodeString(GSHEET_SECRET_BASE64)
	if err != nil {
		log.Fatal("Base64 Decoding Error:", err)
	}

	conf, err := google.JWTConfigFromJSON(GSHEET_SECRET, spreadsheet.Scope)
	if err != nil {
		log.Fatal("Connection Error:", err)
	}

	client := conf.Client(context.TODO())
	Service = spreadsheet.NewServiceWithClient(client)
}
