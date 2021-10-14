package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now()

	// Création d'une image HD pour tester
	img := slice2D(1920, 1080)
	for i := 0; i < len(img); i++ {
		for j := 0; j < len(img[0]); j++ {
			img[i][j] = int16(j + 10*i)
		}
	}
	// Création du kernel identité pour tester
	kernel := [][]int16{
		{0, -1, 0},
		{-1, 4, -1},
		{0, -1, 0},
	}

	final := convolute(img, kernel)
	elapsed := time.Since(start)
	fmt.Printf("Temps: %s\n", elapsed)
	//fmt.Printf("%v", final)
	_ = final
}
