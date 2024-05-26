package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/xManan/fusion-cart/internal/app"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	mux := http.NewServeMux()

	mongoUri := os.Getenv("MONGODB_URI")
	dbName := os.Getenv("MONGODB_DB_NAME")

	err = app.MongoInit(mongoUri, dbName)
	if err != nil {
		log.Fatal(err)
	}

	defer func ()  {
		if err := app.MongoClose(); err != nil {
			log.Fatal(err)
		}
	}()

	app.RegisterRoutes(mux)

	fmt.Printf("Listening at port %s ...", port)
	http.ListenAndServe(":" + port, mux)
}
