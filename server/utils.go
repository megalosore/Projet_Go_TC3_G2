package main

import (
	"image"
	"image/color"
)

//Fonctions utilisés dans les convolutions et le serveur

// Ajoute tous les coefficients du kernel
func sum2D(kernel [][]int16) int16 {
	result := int16(0)
	for i := range kernel {
		for j := range kernel[0] {
			result += kernel[i][j]
		}
	}
	return result
}

// Récupérer un carré de l'image originale centré en x,y et de dimension size*size
func crop(imgAgrandie [][]int16, x int, y int, size int) [][]int16 {
	imgResult := slice2D(size, size)

	// On remplit le carré par les valeurs correspondantes
	for ligne := 0; ligne < size; ligne++ {
		for colonne := 0; colonne < size; colonne++ {
			imgResult[ligne][colonne] = imgAgrandie[y+ligne][x+colonne]
		}
	}
	return imgResult
}

// Crée une version entourée de 0 de l'image originale pour traiter les cas des x,y en bordure
func agrandie(image [][]int16) [][]int16 {
	newImage := slice2D(len(image)+2, len(image[0])+2)
	for i := 1; i < len(newImage)-1; i++ {
		for j := 1; j < len(newImage[0])-1; j++ {
			newImage[i][j] = image[i-1][j-1]
		}
	}
	return newImage
}

// Crée un double slice de dimension précisée (y=ligne , x=colonne) rempli de 0
func slice2D(lenY int, lenX int) [][]int16 {
	doubleSlice := make([][]int16, lenY)
	for i := range doubleSlice {
		doubleSlice[i] = make([]int16, lenX)
	}
	return doubleSlice
}

// Convertit une slice 2D en une image en nuances de gris
func sliceToImg(slice [][]int16) *image.Gray {
	img := image.NewGray(image.Rect(0, 0, len(slice[0]), len(slice)))

	for i := 0; i < len(slice); i++ {
		for j := 0; j < len(slice[0]); j++ {
			img.Set(j, i, color.Gray{Y: uint8(slice[i][j])})
		}
	}
	return img
}

// Convertit une image en un tableau 2D d'int16 exploitable pour la convolution, en nuances de gris
func imgToSlice(image image.Image) [][]int16 {
	lenX := image.Bounds().Size().X
	lenY := image.Bounds().Size().Y
	returnImg := slice2D(lenY, lenX)

	for i := 0; i < lenY; i++ {
		for j := 0; j < lenX; j++ {
			// On récupère la couleur du pixel (j, i) qu'on convertit en un struct color.RGBA contenant les valeurs RGB dans la range 0-255
			RGBA := color.RGBAModel.Convert(image.At(j, i)).(color.RGBA)
			returnImg[i][j] = int16((RGBA.R + RGBA.G + RGBA.B) / 3)
		}
	}
	return returnImg
}
