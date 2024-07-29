package main

import (
	imageRest "github.com/devashish94/squeeze-box-backend/internal/image/api/rest"
	"log"
	"net/http"
	"os"
)

func run() error {
	addr := ":4000"
	imageRouter := imageRest.NewRouter()

	router := http.NewServeMux()
	router.Handle("/api/image/", http.StripPrefix("/api/image", imageRouter))

	log.Println("Server has started at PORT:", addr)

	err := http.ListenAndServe(addr, router)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	err := run()
	if err != nil {
		os.Exit(1)
	}
}
