package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/devashish94/squeeze-box-backend/util"
	"github.com/google/uuid"
)

type StandardResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func main() {
	router := http.NewServeMux()

	router.HandleFunc("POST /api/upload-images", handleImageUpload)
	router.HandleFunc("POST /api/download", handleImageDownload)
	router.HandleFunc("/", handleRoot)

	corsEnabledRouter := CorsMiddleware(router)
	log.Fatal(http.ListenAndServe(":4000", corsEnabledRouter))
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Println("this root is hit")
	w.Header().Set("Content-Type", "application/json")
	data, _ := json.Marshal(StandardResponse{Success: true, Message: "This is the default response"})
	_, err := w.Write(data)
	if err != nil {
		log.Fatal(err)
	}
}

func handleImageUpload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("got hit")
	clientID := uuid.New().String()

	fmt.Println(r.Header.Get("Content-Type"))

	w.Header().Set("Content-Type", "application/json")

	file, _, err := r.FormFile("images")
	if err != nil {
		data, _ := json.Marshal(StandardResponse{Success: false, Message: "no formfile images"})
		w.Write(data)
		return
	}
	defer file.Close()

	multipartData := r.MultipartForm
	if multipartData == nil {
		data, _ := json.Marshal(StandardResponse{Success: false, Message: "images not present in the body"})
		w.Write(data)
		return
	}

	files := multipartData.File["images"]
	targetSize, _ := strconv.ParseFloat(r.FormValue("targetSize"), 64)

	fmt.Println("main.go size ->", targetSize)

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			data, _ := json.Marshal(StandardResponse{Success: false, Message: "could not open fileHeader for " + fileHeader.Filename})
			w.Write(data)
			return
		}
		defer file.Close()

		uniqueDirectoryName := "./uploads/" + clientID
		os.MkdirAll(uniqueDirectoryName, os.ModePerm)

		filePath := filepath.Join("./uploads/"+clientID, fileHeader.Filename)
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

	util.CompressImage(clientID, targetSize)

	data, _ := json.Marshal(StandardResponse{Success: true, Message: clientID})
	w.Write(data)
}

func handleImageDownload(w http.ResponseWriter, r *http.Request) {
	var requestBody map[string]string

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		data, _ := json.Marshal(StandardResponse{Success: false, Message: "could not read the client id"})
		w.Write(data)
		return
	}

	clientID := requestBody["clientID"]
	imagesDirectory := filepath.Join("./output-images", clientID)
	tempZipLocation := filepath.Join("./output-images", clientID, "temp.zip")

	cmd := exec.Command("zip", "-jr", tempZipLocation, imagesDirectory)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		data, _ := json.Marshal(StandardResponse{Success: false, Message: "could not read the client id"})
		w.Write(data)
		panic(err)
	}

	zipFile, err := os.Open(tempZipLocation)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		data, _ := json.Marshal(StandardResponse{Success: false, Message: "could not open the temp.zip file"})
		w.Write(data)
		panic(err)
	}
	defer zipFile.Close()

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=temp.zip")

	_, err = io.Copy(w, zipFile)
	if err != nil {
		log.Fatal(err)
	}
}
