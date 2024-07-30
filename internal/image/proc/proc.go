package proc

import (
	"fmt"
	"github.com/google/uuid"
	_ "image/gif"
	_ "image/png"
	"mime/multipart"
	"os"
)

func CompressImagesToTargetSize(fileHeaders []*multipart.FileHeader, targetSize int64) (string, error) {
	clientID, err := createOutputImagesDirectory()
	if err != nil {
		return "", err
	}

	err = compressImages(fileHeaders, clientID, targetSize)
	if err != nil {
		return "", err
	}

	return clientID, nil
}

func createOutputImagesDirectory() (string, error) {
	clientID := uuid.New().String()
	outputImagesPath := fmt.Sprintf("output-images/%s/", clientID)

	err := os.MkdirAll(outputImagesPath, os.ModePerm)
	if err != nil {
		return "", err
	}

	return clientID, nil
}
