package proc

import (
	"image"
	"image/jpeg"
	"log"
	"os"
)

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
