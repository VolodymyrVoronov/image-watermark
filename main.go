package main

import (
	"fmt"
	"image-watermark/utils"
	"os"
	"time"
)

const inputDirPath = "./input"

// const outputDirPath = "./output"

func main() {
	start := time.Now()
	watermark := utils.GetWatermark()

	inputDir, err := os.ReadDir(inputDirPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	doneChannels := make([]chan bool, len(inputDir))
	errorChannels := make([]chan error, len(inputDir))

	for i, entry := range inputDir {

		if !entry.IsDir() {
			entryName := entry.Name()

			doneChannel := make(chan bool, 1)
			errorChannel := make(chan error, 1)
			doneChannels[i] = doneChannel
			errorChannels[i] = errorChannel

			go utils.ProcessImage(inputDirPath, entryName, watermark, doneChannel, errorChannel)
		}
	}

	for i := range inputDir {
		select {
		case done := <-doneChannels[i]:
			if done {
				fmt.Println("\x1b[1mImage\x1b[0m \x1b[34m" + inputDir[i].Name() + "\x1b[0m was processed \x1b[32msuccessfully!\x1b[0m")
			}

		case err := <-errorChannels[i]:
			if err != nil {
				fmt.Println("\x1b[1mImage\x1b[0m \x1b[34m" + inputDir[i].Name() + "\x1b[0m was processed \x1b[31mwith error!\x1b[0m")
				fmt.Println(err)
			}
		}
	}

	duration := time.Since(start)

	fmt.Println()
	fmt.Println("Done!")
	fmt.Println("Total time: ", duration)
}
