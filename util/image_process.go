package util

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/chai2010/webp"
)

func CompressImage(dir *os.File, limit int) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	filenames, err := dir.Readdirnames(-1)
	HandleError(err)

	// limit := int64(300)
	targetSize := int64(limit * 1000)
	wg := sync.WaitGroup{}

	wg.Add(len(filenames))
	parallel := time.Now()
	for _, filename := range filenames {
		go func(filename string) {
			ProcessImage(filename, targetSize)
			defer wg.Done()
		}(filename)
	}
	wg.Wait()
	fmt.Println("Time taken (parallel):", time.Since(parallel))
}

func ProcessImage(filename string, targetSize int64) {
	filenameWithoutExtension := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
	outputImagePath := "./output-images/" + filenameWithoutExtension + ".jpeg"
	inputImagePath := "./input-images/" + filename

	file, err := os.Open(inputImagePath)
	HandleError(err)
	defer file.Close()

	fileInfo, err := os.Stat(inputImagePath)
	HandleError(err)
	if fileInfo.Size() <= targetSize {
		cmd := exec.Command("cp", inputImagePath, outputImagePath)
		cmd.Run()
		fmt.Println("[", filename, "] Size:", float64(fileInfo.Size())/float64(1000), "KB EARLY RETURN")
		return
	}

	img, err := DecodeImage(file, filename)
	HandleError(err)

	currentSize := int64(math.MaxInt64)
	qualityFactor := 100

	left := 0
	right := qualityFactor
	var qualityFactorAnswer int

	for left <= right {
		mid := left + (right-left)/2

		outputFile, err := os.Create(outputImagePath)
		HandleError(err)
		defer outputFile.Close()

		err = jpeg.Encode(outputFile, img, &jpeg.Options{Quality: mid})
		HandleError(err)

		info, err := os.Stat(outputImagePath)
		HandleError(err)
		currentSize = info.Size()

		if currentSize <= targetSize {
			qualityFactorAnswer = mid
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	outputFile, err := os.Create(outputImagePath)
	HandleError(err)
	defer outputFile.Close()

	err = jpeg.Encode(outputFile, img, &jpeg.Options{Quality: qualityFactorAnswer})
	HandleError(err)

	info, err := os.Stat(outputImagePath)
	HandleError(err)
	fmt.Println("[", filenameWithoutExtension+".jpeg", "] Size:", float64(info.Size())/float64(1000), "KB")
}

func DecodeImage(file *os.File, filename string) (image.Image, error) {
	ext := filepath.Ext(filename)

	switch ext {
	case ".jpg", ".jpeg":
		return jpeg.Decode(file)
	case ".png":
		return png.Decode(file)
	case ".webp":
		return webp.Decode(file)
	default:
		return nil, fmt.Errorf("file Type not supported, add it in the decodeImage() %s", ext)
	}
}

// LINEAR SEARCH IMPLEMENTATION
// for qualityFactor > 0 && currentSize > targetSize {
// 	outputFile, err := os.Create(outputImagePath)
// 	HandleError(err)
// 	err = jpeg.Encode(outputFile, img, &jpeg.Options{Quality: qualityFactor})
// 	outputFile.Close()
// 	infoo, err := os.Stat(outputImagePath)
// 	HandleError(err)
// 	currentSize = infoo.Size()
// 	// don't know how to fine tune this
// 	qualityFactor -= qualityFactorCalculate(currentSize, targetSize)
// }

// dir, err := os.Open("./input-images")
// util.HandleError(err)
// defer dir.Close()

// util.CompressImage(dir, 200)
