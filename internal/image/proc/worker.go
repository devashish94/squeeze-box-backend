package proc

import (
	"log"
	"mime/multipart"
	"runtime"
	"sync"
)

type Job struct {
	fileHeader *multipart.FileHeader
	clientID   string
	targetSize int64
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
