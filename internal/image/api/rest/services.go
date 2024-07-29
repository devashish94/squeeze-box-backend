package rest

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"os/exec"
	"strconv"
)

func ParseMultipartForm(r *http.Request) ([]*multipart.FileHeader, int64, error) {
	err := r.ParseMultipartForm(10_000_000)
	if err != nil {
		return nil, 0, err
	}

	multipartData := r.MultipartForm
	if multipartData == nil {
		return nil, 0, errors.New("multipart form data is nil")
	}

	files := multipartData.File["images"]
	targetSize, err := strconv.ParseFloat(r.FormValue("targetSize"), 64)
	if err != nil {
		return nil, 0, err
	}

	return files, int64(targetSize * 1000), nil
}

func CreateZip(clientID string) (string, error) {
	outputImagesDirectory := fmt.Sprintf("output-images/%s", clientID)
	zipFilePath := fmt.Sprintf("%s/temp.zip", outputImagesDirectory)

	cmd := exec.Command("zip", "-jr", zipFilePath, outputImagesDirectory)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return zipFilePath, nil
}
