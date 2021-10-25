package main

import "sync"

// Effectue la convolution pour 1 pixel
func computeconvolution(resultArray [][]int16, imgAgrandie [][]int16, kernel [][]int16, seuil float64, x int, y int) {
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
		resultArray[y][x] = int16(result / somme)
	} else {
		//seuillage pour les detections de contours
		if result < int16(255*seuil) {
			resultArray[y][x] = 0
		} else {
			resultArray[y][x] = 255
		}
	}
}

// Calcule la convolution sur un certain nombre de lignes
func lineCompute(resultArray [][]int16, imageAgrandie [][]int16, kernel [][]int16, seuil float64, y int, lenX int, nbLigne int, waitGroup *sync.WaitGroup) {
	for i := y; i < y+nbLigne; i++ {
		for j := 0; j < lenX; j++ {
			// Calcule la convolution pour un pixel
			computeconvolution(resultArray, imageAgrandie, kernel, seuil, j, i)
		}
	}
	waitGroup.Done()
}

// Fonction à appeler pour effectuer la convolution d'une image et d'un filtre
func convolute(imageArray [][]int16, kernel [][]int16, seuil float64) [][]int16 {
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
			go lineCompute(result, imageAgrandie, kernel, seuil, i, lenX, lenY-i, &waitGroup)
		} else {
			go lineCompute(result, imageAgrandie, kernel, seuil, i, lenX, nbLigne, &waitGroup)
		}
	}
	waitGroup.Wait()
	return result
}
