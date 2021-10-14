package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	img, err := openImg("")
	if err != nil {
		panic(err)
	}
	imgConverted := imgToSlice(img)
	kernel := [][]int16{
		{0, -1, 0},
		{-1, 4, -1},
		{0, -1, 0},
	}

	final := convolute(imgConverted, kernel)
	finalImage := sliceToImg(final)
	elapsed := time.Since(start)
	fmt.Printf("Temps: %s\n", elapsed)
	writeImg(finalImage, "output")
}
