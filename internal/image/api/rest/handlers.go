package rest

import (
	"github.com/devashish94/squeeze-box-backend/internal/image/proc"
	"github.com/devashish94/squeeze-box-backend/internal/response"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func ImageUploadHandler(w http.ResponseWriter, r *http.Request) {
	res := response.New(w)

	imageFiles, targetSize, err := ParseMultipartForm(r)
	if err != nil {
		log.Println("Problem parsing multipart form:", err.Error())
		res.Status(400).Json(response.JsonResponse{
			Status:     "fail",
			StatusCode: http.StatusBadRequest,
			Data:       err.Error(),
		})
		return
	}

	startTime := time.Now()
	clientID, err := proc.CompressImagesToTargetSize(imageFiles, targetSize)
	if err != nil {
		log.Println("Problem compressing images:", err.Error())
		res.Status(500).Json(response.JsonResponse{
			Status:     "fail",
			StatusCode: http.StatusInternalServerError,
			Data:       err.Error(),
		})
		return
	}
	log.Println("Total Time:", time.Since(startTime), clientID)

	res.Json(response.JsonResponse{
		Status:     "ok",
		StatusCode: http.StatusOK,
		Data: map[string]interface{}{
			"clientId": clientID,
		},
	})
}

func ImageDownloadHandler(w http.ResponseWriter, r *http.Request) {
	clientID := r.PathValue("clientId")
	res := response.New(w)

	zipFilePath, err := CreateZip(clientID)
	if err != nil {
		res.Status(500).Json(response.JsonResponse{
			Status:     "fail",
			StatusCode: http.StatusInternalServerError,
			Data:       "Could not zip the images",
		})
		return
	}

	zipFile, err := os.Open(zipFilePath)
	if err != nil {
		res.Status(500).Json(response.JsonResponse{
			Status:     "fail",
			StatusCode: http.StatusInternalServerError,
			Data:       "Could not open the zip file",
		})
		return
	}
	defer func(zipFile *os.File) {
		err := zipFile.Close()
		if err != nil {
			log.Println("Could not close the zip file")
		}
	}(zipFile)

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=shrunk-images.zip")

	if _, err := io.Copy(w, zipFile); err != nil {
		res.Status(500).Json(response.JsonResponse{
			Status:     "fail",
			StatusCode: http.StatusInternalServerError,
			Data:       "Could not write the zip file to response",
		})
	}
}
