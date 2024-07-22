package main

import (
	"fmt"
	"github.com/devashish94/squeeze-box-backend/internal/api/rest"
	"net/http"
	"os"
)

func run() error {
	addr := ":4000"
	apiRouter := rest.NewRouter()

	fmt.Println("Server has started at PORT:", addr)

	err := http.ListenAndServe(addr, apiRouter)
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
