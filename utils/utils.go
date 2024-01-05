package utils

import (
	"bufio"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

func GetWatermark() image.Image {
	watermarkDir, err := os.ReadDir("./watermark")
	if err != nil {
		fmt.Println(err)
		return nil
	}

	if len(watermarkDir) == 0 {
		fmt.Println("Watermark not found.")
		return nil
	}

	if len(watermarkDir) > 1 {
		fmt.Println("More than one watermark found.")
		fmt.Println("Please provide only one watermark.")
		return nil
	}

	watermark, err := os.Open("./watermark/" + watermarkDir[0].Name())
	if err != nil {
		fmt.Println(err)
		return nil
	}

	defer watermark.Close()

	decodedWatermark, err := png.Decode(watermark)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return decodedWatermark
}

func GetImageFormat(file string) (string, error) {
	extension := filepath.Ext(file)

	if extension == "" {
		return "", fmt.Errorf("unable to determine image format: file has no extension")
	}

	format := strings.ToLower(extension[1:])

	return format, nil
}

func ProcessImage(inputDirPath string, entryName string, watermark image.Image, doneChannel chan bool, errorChannel chan error) {
	file, err := os.Open(inputDirPath + "/" + entryName)
	if err != nil {
		fmt.Println(err)
		errorChannel <- err
		return
	}
	defer file.Close()

	format, err := GetImageFormat(entryName)
	if err != nil {
		fmt.Println(err)
		errorChannel <- err
		return
	}

	var decodedImage image.Image

	if format == "png" {
		decodedImage, err = png.Decode(file)
		if err != nil {
			fmt.Println("Error while decoding png: ", err)
			errorChannel <- err
			return
		}
	}

	if format == "jpeg" {
		decodedImage, err = jpeg.Decode(file)
		if err != nil {
			fmt.Println("Error while decoding jpeg: ", err)
			errorChannel <- err
			return
		}
	}

	b := decodedImage.Bounds()
	m := image.NewRGBA(b)
	draw.Draw(m, b, decodedImage, image.Point{}, draw.Src)

	horizontalSpacing := 50
	verticalSpacing := 50

	for y := 0; y < b.Dy(); y += watermark.Bounds().Dy() + verticalSpacing {
		for x := 0; x < b.Dx(); x += watermark.Bounds().Dx() + horizontalSpacing {
			offset := image.Pt(x, y)
			draw.Draw(m, watermark.Bounds().Add(offset), watermark, image.Point{}, draw.Over)
		}
	}

	var result *os.File

	if format == "png" {
		result, err = os.Create(fmt.Sprintf("./output/%s", entryName))
		if err != nil {
			fmt.Println(err)
			errorChannel <- err
			return
		}

		png.Encode(result, m)
	}

	if format == "jpeg" {
		result, err = os.Create(fmt.Sprintf("./output/%s", entryName))
		if err != nil {
			fmt.Println(err)
			errorChannel <- err
			return
		}

		jpeg.Encode(result, m, &jpeg.Options{Quality: jpeg.DefaultQuality})
	}

	doneChannel <- true

	defer result.Close()
}

func ClearInputDir(inputDirPath string) error {
	err := filepath.Walk(inputDirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			err = os.Remove(path)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		fmt.Println(err)
	}

	return nil
}

func GetUserInputForClearInputDir(inputDirPath string) error {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println()
	fmt.Print("Clear input directory? (y/n): ")
	scanner.Scan()
	userInput := scanner.Text()

	if userInput == "y" || userInput == "Y" || userInput == "yes" || userInput == "Yes" {
		err := ClearInputDir(inputDirPath)
		if err != nil {
			fmt.Println(err)
			return err
		}

		fmt.Println()
		fmt.Println("Input directory was cleared!")
		fmt.Println("Exiting...")
	} else if userInput == "n" || userInput == "N" || userInput == "no" || userInput == "No" {
		fmt.Println()
		fmt.Println("Input directory was not cleared!")
		fmt.Println("Exiting...")
	} else {
		fmt.Println()
		fmt.Println("Invalid input!")
		fmt.Println("Exiting...")
	}

	return nil
}
