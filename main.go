package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

// "github.com/devashish94/pixel-presser/util"

type StandardResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func main() {
	http.HandleFunc("POST /api/upload-images", handleImageUpload)
	http.HandleFunc("/", handleRoot)

	err := http.ListenAndServe(":6900", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data, _ := json.Marshal(StandardResponse{Success: true, Message: "This is the default response"})
	_, err := w.Write(data)
	if err != nil {
		log.Fatal(err)
	}
}

func handleImageUpload(w http.ResponseWriter, r *http.Request) {
	clientID := uuid.New().String()

	multipartData := r.MultipartForm
	files := multipartData.File["images"]

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			data, _ := json.Marshal(StandardResponse{Success: false, Message: "could not read multipart form"})
			w.Write(data)
			return
		}
		defer file.Close()

		uniqueDirectoryName := "./uploads" + clientID
		os.MkdirAll(uniqueDirectoryName, os.ModePerm)

		filePath := filepath.Join("./uploads", fileHeader.Filename)
		newFile, err := os.Create(filePath)
		if err != nil {
			data, _ := json.Marshal(StandardResponse{Success: false, Message: "could not create file" + filePath})
			w.Write(data)
			return
		}
		defer newFile.Close()

		io.Copy(newFile, file)

		fmt.Println("done uploading", filePath)
	}

	data, _ := json.Marshal(StandardResponse{Success: true, Message: clientID})
	w.Write(data)
}
