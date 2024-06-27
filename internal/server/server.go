package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"web-storage-service/pkg"

	"web-storage-service/internal/database"
)

var (
	port int
)

func init() {
	// TODO: Если использовать "github.com/joho/godotenv/autoload",
	// TODO: можно будет убрать init() функции
	err := pkg.LoadEnv(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	port, _ = strconv.Atoi(os.Getenv("PORT"))
}

type Server struct {
	port int

	db database.Service
}

func NewServer() *http.Server {
	NewServer := &Server{
		port: port,

		db: database.New(),
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
