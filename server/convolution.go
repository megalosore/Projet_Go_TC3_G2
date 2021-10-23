package main

import (
	"image"
	"image/color"
	"sync"
)

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

// Crée un double slice de dimension précisée (y=ligne , x=colonne) rempli de 0
func slice2D(lenY int, lenX int) [][]int16 {
	doubleSlice := make([][]int16, lenY)
	for i := range doubleSlice {
		doubleSlice[i] = make([]int16, lenX)
	}
	return doubleSlice
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

// Récupére un carré de l'image originale centré en x,y et de dimension size*size
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

// Effectue la convolution pour 1 pixel
func computeconvolution(resultArray [][]int16, imgAgrandie [][]int16, kernel [][]int16, x int, y int) {
	size := len(kernel)
	// On récupère une version size*size entourant le pixel que l'on veut traiter
	cropedImage := crop(imgAgrandie, x, y, size)
	result := int16(0)
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			// On effectue le calcul de la convolution on ajoutant les éléments opposés entre le filtre et l'image cropée
			result += kernel[i][j] * cropedImage[size-i-1][size-j-1]
		}
	}
	if result < 0 { //On borne notre resultat entre 0 et 255
		result = -1 * result
	}
	if result > 255 {
		result = 255
	}

	// On normalise le résultat par la somme des coefficients du filtre si le filtre le permet
	somme := sum2D(kernel)
	if somme != 0 {
		resultArray[y][x] = result / somme
	} else {
		resultArray[y][x] = result
	}
}

// Calcule la convolution sur un certain nombre de lignes
func lineCompute(resultArray [][]int16, imageAgrandie [][]int16, kernel [][]int16, y int, lenX int, nbLigne int, waitGroup *sync.WaitGroup) {
	for i := y; i < y+nbLigne; i++ {
		for j := 0; j < lenX; j++ {
			// Calcule la convolution pour un pixel
			computeconvolution(resultArray, imageAgrandie, kernel, j, i)
		}
	}
	waitGroup.Done()
}

// Fonction à appeler pour effectuer la convolution d'une image et d'un filtre
func convolute(imageArray [][]int16, kernel [][]int16) [][]int16 {
	lenY := len(imageArray)
	lenX := len(imageArray[0])
	// On traite l'image pour rajouter des 0 sur les bordures
	imageAgrandie := agrandie(imageArray)
	result := slice2D(lenY, lenX)
	// On définit le nombre de go routine max
	const nbRoutine = 12
	// On rajoute 1 pour éviter les cas ou nbLigne est arrondi à l'inférieur
	nbLigne := (lenY / nbRoutine) + 1
	var waitGroup sync.WaitGroup

	for i := 0; i < lenY; i += nbLigne {
		waitGroup.Add(1)
		// On vérifie qu'on ne dépasse pas le nombre de lignes max
		if i+nbLigne > lenY {
			go lineCompute(result, imageAgrandie, kernel, i, lenX, lenY-i, &waitGroup)
		} else {
			go lineCompute(result, imageAgrandie, kernel, i, lenX, nbLigne, &waitGroup)
		}
	}
	waitGroup.Wait()
	return result
}
