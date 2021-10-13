package main

import (
	"fmt"
	"time"
)

func slice2D(lenY int, lenX int) [][]int16 { // Crée un double slice de dimension précisée (y=ligne , x=colonne) rempli de 0
	doubleSlice := make([][]int16, lenY)
	for i := range doubleSlice {
		doubleSlice[i] = make([]int16, lenX)
	}
	return doubleSlice
}

func agrandie(image [][]int16) [][]int16 {
	newImage := slice2D(len(image)+2, len(image[0])+2) // On recrée une version entourée de 0 de l'image originale pour traiter les cas des x,y en bordures
	for i := 1; i < len(newImage)-1; i++ {
		for j := 1; j < len(newImage[0])-1; j++ {
			newImage[i][j] = image[i-1][j-1]
		}
	}
	return newImage
}

func crop(imgAgrandie [][]int16, x int, y int, size int) [][]int16 { // Récupére un carré de l'image originale centré en x,y et de dimension size*size
	imgResult := slice2D(size, size)

	for ligne := 0; ligne < size; ligne++ { // On remplit le carré par les valeurs correspondantes
		for colonne := 0; colonne < size; colonne++ {
			imgResult[ligne][colonne] = imgAgrandie[y+ligne][x+colonne]
		}
	}
	return imgResult
}

func sum2D(kernel [][]int16) int16 { // Ajoute tous les coefficients du kernel
	result := int16(0)
	for i := range kernel {
		for j := range kernel[0] {
			result += kernel[i][j]
		}
	}
	return result
}

func computeconvolution(resultArray [][]int16, imgAgrandie [][]int16, kernel [][]int16, x int, y int) {
	size := len(kernel)
	cropedImage := crop(imgAgrandie, x, y, size) // On récupère une version size*size entourant le pixel que l'on veut traiter
	result := int16(0)
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			result += kernel[i][j] * cropedImage[size-i-1][size-j-1] // On effectue le calcul de la convolution on ajoutant les éléments opposés entre le filtre et l'image cropé
		}
	}
	somme := sum2D(kernel) // On normalise le résultat par la somme des coefficients du filtre si le filtre le permet
	if somme != 0 {
		resultArray[y][x] = result / somme
	} else {
		resultArray[y][x] = result
	}
}

func lineCompute(resultArray [][]int16, imageAgrandie [][]int16, kernel [][]int16, y int, lenX int, nbLigne int) { // Calcule la convolution sur un certain nombre de lignes
	for i := y; i < y+nbLigne; i++ {
		for j := 0; j < lenX; j++ {
			computeconvolution(resultArray, imageAgrandie, kernel, j, i) // Calcule la convolution pour un pixel
		}
	}
}

func convolute(imageArray [][]int16, kernel [][]int16) [][]int16 { // Fonction à appeler pour effectuer la convolution d'une image et d'un filtre
	lenY := len(imageArray)
	lenX := len(imageArray[0])
	imageAgrandie := agrandie(imageArray) // On traite l'image pour rajouter des 0 sur les bordures
	result := slice2D(lenY, lenX)
	const nbRoutine = 12              // On définit le nombre de go routine max
	nbLigne := (lenY / nbRoutine) + 1 // On rajoute 1 pour éviter les cas ou lenY proche de 12

	for i := 0; i < lenY; i += nbLigne {
		if i+nbLigne > lenY { // On vérifie qu'on ne dépasse par le nombre de lignes max
			go lineCompute(result, imageAgrandie, kernel, i, lenX, lenY-i)
		} else {
			go lineCompute(result, imageAgrandie, kernel, i, lenX, nbLigne)
		}
	}
	return result
}

func main() {
	start := time.Now()

	image := slice2D(1920, 1080) // Création d'une image HD pour tester
	for i := 0; i < len(image); i++ {
		for j := 0; j < len(image[0]); j++ {
			image[i][j] = int16(j + 10*i)
		}
	}
	kernel := slice2D(3, 3) // Création du kernel identité pour tester
	kernel[1][1] = 1

	final := convolute(image, kernel)
	elapsed := time.Since(start)
	fmt.Printf("Temps: %s\n", elapsed)
	_ = final
}
