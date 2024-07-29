package rest

import "net/http"

func NewRouter() http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("POST /upload", ImageUploadHandler)
	router.HandleFunc("GET /download/{clientId}", ImageDownloadHandler)
	return router
}
