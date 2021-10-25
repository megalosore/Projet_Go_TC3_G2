package main

import (
	"math"
	"runtime"
	"sync"
)

// Effectue la convolution pour 1 pixel
func computeConvolutionDouble(resultArray [][]int16, imgAgrandie [][]int16, kernel1 [][]int16, kernel2 [][]int16, seuil float64, x int, y int) {
	size := len(kernel1)
	// On récupère une version size*size entourant le pixel que l'on veut traiter
	cropedImage := crop(imgAgrandie, x, y, size)
	result1 := int16(0)
	result2 := int16(0)
	//var temp1 float64
	//var temp2 float64
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			// On effectue le calcul de la convolution on ajoutant les éléments opposés entre le filtre et l'image cropée
			// On implémente la convolution de Sobel
			//temp1 = math.Pow(float64(kernel1[i][j]*cropedImage[size-i-1][size-j-1]), 2)
			//temp2 = math.Pow(float64(kernel2[i][j]*cropedImage[size-i-1][size-j-1]), 2)
			result1 += kernel1[i][j] * cropedImage[size-i-1][size-j-1]
			result2 += kernel2[i][j] * cropedImage[size-i-1][size-j-1]
		}
	}
	result := math.Sqrt(math.Pow(float64(result1), 2)+math.Pow(float64(result2), 2)) / 4
	if result < seuil*255 { //Seuillage
		resultArray[y][x] = 0
	} else {
		resultArray[y][x] = 255
	}
}

// Calcule la convolution sur un certain nombre de lignes
func lineComputeDouble(resultArray [][]int16, imageAgrandie [][]int16, kernel1 [][]int16, kernel2 [][]int16, seuil float64, y int, lenX int, nbLigne int, waitGroup *sync.WaitGroup) {
	for i := y; i < y+nbLigne; i++ {
		for j := 0; j < lenX; j++ {
			// Calcule la convolution pour un pixel
			computeConvolutionDouble(resultArray, imageAgrandie, kernel1, kernel2, seuil, j, i)
		}
	}
	waitGroup.Done()
}

// Fonction à appeler pour effectuer la convolution d'une image et d'un filtre
func convoluteDouble(imageArray [][]int16, kernel1 [][]int16, kernel2 [][]int16, seuil float64) [][]int16 {
	lenY := len(imageArray)
	lenX := len(imageArray[0])
	// On traite l'image pour rajouter des 0 sur les bordures
	imageAgrandie := agrandie(imageArray)
	result := slice2D(lenY, lenX)
	// On définit le nombre de go routine max
	var nbRoutine = runtime.NumCPU() * 2
	// On rajoute 1 pour éviter les cas ou nbLigne est arrondi à l'inférieur
	nbLigne := (lenY / nbRoutine) + 1
	var waitGroup sync.WaitGroup

	for i := 0; i < lenY; i += nbLigne {
		waitGroup.Add(1)
		// On vérifie qu'on ne dépasse pas le nombre de lignes max
		if i+nbLigne > lenY {
			go lineComputeDouble(result, imageAgrandie, kernel1, kernel2, seuil, i, lenX, lenY-i, &waitGroup)
		} else {
			go lineComputeDouble(result, imageAgrandie, kernel1, kernel2, seuil, i, lenX, nbLigne, &waitGroup)
		}
	}
	waitGroup.Wait()
	return result
}
