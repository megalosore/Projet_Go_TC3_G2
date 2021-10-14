package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	img, _ := openImg("")
	imgConverted := imgToSlice(img)
	kernel := [][]int16{
		{0, -1, 0},
		{-1, 4, -1},
		{0, -1, 0},
	}

	final := convolute(imgConverted, kernel)
	final_image := sliceToImg(final)
	elapsed := time.Since(start)
	fmt.Printf("Temps: %s\n", elapsed)
	writeImg(final_image, "output")
}