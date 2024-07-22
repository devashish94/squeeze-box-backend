package rest

import "net/http"

func NewRouter() http.Handler {
	pingRouter := http.NewServeMux()
	pingRouter.HandleFunc("/ping", PingHandler)

	imageRouter := http.NewServeMux()
	imageRouter.HandleFunc("POST /upload", ImageUploadHandler)
	//imageRouter.HandleFunc("GET /download", nil)

	mainRouter := http.NewServeMux()
	mainRouter.Handle("/ping", pingRouter)
	mainRouter.Handle("/image", imageRouter)
	return mainRouter
}
