package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/akhiltn/go-link/docs"
	"github.com/akhiltn/go-link/internal/app"
	"github.com/akhiltn/go-link/internal/data"
	"github.com/joho/godotenv"
)

// @title Go Link
// @version 1.0
// @description URL shortener app
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /
func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	db, err := data.New(ctx, data.Valkey, os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}
	defer db.Close()

	server := app.New(db, app.Config{Domain: os.Getenv("DOMAIN")})

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000" // Default port if not set
	}
	log.Printf("Server started on port %s", port)
	if err := server.Listen(":" + port); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	return nil
}
