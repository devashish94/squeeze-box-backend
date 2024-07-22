package rest

import (
	"fmt"
	"github.com/devashish94/squeeze-box-backend/internal/response"
	"github.com/devashish94/squeeze-box-backend/util"
	"github.com/google/uuid"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func PingHandler(w http.ResponseWriter, r *http.Request) {
	res := response.New(w)
	res.Json(JsonResponse{
		Status:     "ok",
		StatusCode: http.StatusOK,
		Data:       "Server is running normally",
	})
}

func ImageUploadHandler(w http.ResponseWriter, r *http.Request) {
	clientID := uuid.New().String()
	res := response.New(w)

	multipartData := r.MultipartForm
	if multipartData == nil {
		res.Json(JsonResponse{
			Status:     "ok",
			StatusCode: http.StatusBadRequest,
			Data:       "Multipart form data is null",
		})
		return
	}

	files := multipartData.File["images"]
	targetSize, _ := strconv.ParseFloat(r.FormValue("targetSize"), 64)

	fmt.Println("targetSize from client:", targetSize)

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			res.Json(JsonResponse{
				Status:     "ok",
				StatusCode: http.StatusBadRequest,
				Data:       "Could not open file header for " + fileHeader.Filename,
			})
			return
		}

		uniqueDirectoryName := "./uploads/" + clientID
		_ = os.MkdirAll(uniqueDirectoryName, os.ModePerm)

		filePath := filepath.Join("./uploads/"+clientID, fileHeader.Filename)
		newFile, err := os.Create(filePath)
		if err != nil {
			res.Json(JsonResponse{
				Status:     "ok",
				StatusCode: http.StatusInternalServerError,
				Data:       "Could not create file for output image: " + fileHeader.Filename,
			})
			return
		}
		_ = file.Close()
		_ = newFile.Close()

		_, _ = io.Copy(newFile, file)

		fmt.Println("done uploading", filePath)
	}

	util.CompressImage(clientID, targetSize)

	res.Json(JsonResponse{
		Status:     "ok",
		StatusCode: http.StatusOK,
		Data: map[string]interface{}{
			"clientId": clientID,
		},
	})
}
