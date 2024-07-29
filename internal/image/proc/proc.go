package proc

import (
	"fmt"
	"github.com/google/uuid"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"log"
	"mime/multipart"
	"os"
	"runtime"
	"sync"
)

type Job struct {
	fileHeader *multipart.FileHeader
	clientID   string
	targetSize int64
}

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

func worker(jobChannel <-chan *Job, errChan chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobChannel {
		err := compressImage(job.fileHeader, job.clientID, job.targetSize)
		log.Println(job.clientID, job.fileHeader.Filename)
		if err != nil {
			errChan <- err
		}
	}
}

func startWorkers(jobChan chan *Job, errChan chan error, waitGroup *sync.WaitGroup) {
	numOfWorkers := runtime.NumCPU()
	for core := 0; core < numOfWorkers; core++ {
		waitGroup.Add(1)
		go worker(jobChan, errChan, waitGroup)
	}
}

func pushJobs(jobChan chan *Job, fileHeaders []*multipart.FileHeader, clientID string, targetSize int64) {
	for _, fileHeader := range fileHeaders {
		jobChan <- &Job{fileHeader: fileHeader, clientID: clientID, targetSize: targetSize}
	}
}

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

func calculateOptimalQualityFactor(img *image.Image, filepath string, currentFileSize, targetSize int64) (int, error) {
	lowerQualityFactor := 0
	higherQualityFactor := 100
	optimalQualityFactor := 101

	for lowerQualityFactor <= higherQualityFactor {
		midQualityFactor := lowerQualityFactor + (higherQualityFactor-lowerQualityFactor)/2
		currentFileSize, err := calculateNewFileSize(midQualityFactor, filepath, img)
		if err != nil {
			log.Println("Could not calculate new file size", err.Error())
			return 0, err
		}

		if currentFileSize > targetSize {
			optimalQualityFactor = midQualityFactor
			higherQualityFactor = midQualityFactor - 1
		} else {
			lowerQualityFactor = midQualityFactor + 1
		}
	}

	if optimalQualityFactor > 0 && currentFileSize > targetSize {
		optimalQualityFactor--
	}

	return optimalQualityFactor, nil
}

func calculateNewFileSize(qualityFactor int, filepath string, image *image.Image) (int64, error) {
	file, err := os.Create(filepath)
	if err != nil {
		return 0, err
	}

	err = jpeg.Encode(file, *image, &jpeg.Options{Quality: qualityFactor})
	if err != nil {
		return 0, err
	}

	fileStat, err := file.Stat()
	if err != nil {
		return 0, err
	}
	return fileStat.Size(), nil
}

func encodeImageToQualityFactor(img *image.Image, filepath string, qualityFactor int) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println("Could not close file after creating:", err.Error())
		}
	}(file)

	err = jpeg.Encode(file, *img, &jpeg.Options{Quality: qualityFactor})
	if err != nil {
		return err
	}

	return nil
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
