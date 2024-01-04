package main

import (
	"fmt"
	"image-watermark/utils"
	"os"
	"time"
)

const inputDirPath = "./input"

func main() {
	start := time.Now()
	watermark := utils.GetWatermark()

	inputDir, err := os.ReadDir(inputDirPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(inputDir) == 0 {
		fmt.Println("Input directory is empty!")
		fmt.Println("Exiting...")
		return
	}

	doneChannels := make([]chan bool, len(inputDir))
	errorChannels := make([]chan error, len(inputDir))

	for i, entry := range inputDir {
		doneChannels[i] = make(chan bool, 1)
		errorChannels[i] = make(chan error, 1)

		if !entry.IsDir() {
			entryName := entry.Name()

			go utils.ProcessImage(inputDirPath, entryName, watermark, doneChannels[i], errorChannels[i])
		}
	}

	for i := range inputDir {
		select {
		case done := <-doneChannels[i]:
			if done {
				fmt.Println("File " + inputDir[i].Name() + " was processed successfully!")
			}

		case err := <-errorChannels[i]:
			if err != nil {
				fmt.Println("Image" + inputDir[i].Name() + " was processed with error!")
				fmt.Println(err)
			}
		}
	}

	duration := time.Since(start)

	fmt.Println()
	fmt.Println("Done!")
	fmt.Println("Total processing time: ", duration)

	err = utils.GetUserInputForClearInputDir(inputDirPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	time.Sleep(time.Second)
}
