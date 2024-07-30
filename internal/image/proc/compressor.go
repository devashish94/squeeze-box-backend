package proc

import (
	"fmt"
	"image"
	"log"
	"mime/multipart"
	"sync"
)

func compressImages(fileHeaders []*multipart.FileHeader, clientID string, targetSize int64) error {
	var waitGroup sync.WaitGroup
	errChan := make(chan error)
	jobChan := make(chan *Job, len(fileHeaders))

	startWorkers(jobChan, errChan, &waitGroup)
	pushJobs(jobChan, fileHeaders, clientID, targetSize)

	close(jobChan)
	waitGroup.Wait()
	close(errChan)

	if len(errChan) > 0 {
		log.Println("Problem compressing images in clientID:", clientID)
		return <-errChan
	}

	return nil
}

func compressImage(fileHeader *multipart.FileHeader, clientID string, targetSize int64) error {
	if fileHeader.Size <= targetSize {
		log.Println("Early return, image size already less than target size")
		return nil
	}

	file, err := fileHeader.Open()
	if err != nil {
		log.Println("Could not extract file from the fileHeader", fileHeader.Filename)
		return err
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			log.Println("Could not close the multipart file:", fileHeader.Filename)
		}
	}(file)

	img, _, err := image.Decode(file)
	if err != nil {
		log.Println("Could not decode the image:", fileHeader.Filename)
		return err
	}

	outputFilePath := fmt.Sprintf("output-images/%s/%s", clientID, fileHeader.Filename)
	optimalQualityFactor, err := calculateOptimalQualityFactor(&img, outputFilePath, fileHeader.Size, targetSize)
	if err != nil {
		log.Println("Could not calculate the optimal quality factor", fileHeader.Size, fileHeader.Filename)
		return err
	}

	err = encodeImageToQualityFactor(&img, outputFilePath, optimalQualityFactor)
	if err != nil {
		log.Println("Could not encode the image with the optimal quality factor", fileHeader.Filename)
		return err
	}

	return nil
}
