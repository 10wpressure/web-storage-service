package main

import (
	"fmt"
	"log"
	"os"
	"web-storage-service/internal/server"
	"web-storage-service/pkg"
)

var (
	certFile string
	keyFile  string
)

func init() {
	// TODO: Если использовать "github.com/joho/godotenv/autoload",
	// TODO: можно будет убрать init() функции
	err := pkg.LoadEnv(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	certFile = os.Getenv("CERT_FILE")
	keyFile = os.Getenv("KEY_FILE")
}

func main() {
	srv := server.NewServer()

	certAndKeyExist := pkg.DoFilesExist(certFile, keyFile)

	if !certAndKeyExist {
		err := pkg.GenerateCertificate(certFile, keyFile)
		if err != nil {
			panic(fmt.Sprintf("cannot generate certificate: %s", err))
		}
	}

	certAndKeyExist = pkg.DoFilesExist(certFile, keyFile)

	if certAndKeyExist {
		err := srv.ListenAndServeTLS(certFile, keyFile)
		if err != nil {
			panic(fmt.Sprintf("cannot start server: %s", err))
		}
	} else {
		err := srv.ListenAndServe()
		if err != nil {
			panic(fmt.Sprintf("cannot start server: %s", err))
		}
	}
}
