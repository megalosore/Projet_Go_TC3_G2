package main

import (
	"fmt"
	"time"
)

func slice2D(leny int, lenx int) [][]int { //Crée un double slice de dimenssion précisé(y=ligne , x=collone) rempli de 0
	double_slice := make([][]int, leny)
	for i := range double_slice {
		double_slice[i] = make([]int, lenx)
	}
	return double_slice
}

func agrandit(image [][]int) [][]int {
	new_image := slice2D(len(image)+2, len(image[0])+2) //On recrée une version entouré de 0 de l'image originale pour traiter les cas des x,y en bordures
	for i := 1; i < len(new_image)-1; i++ {
		for j := 1; j < len(new_image[0])-1; j++ {
			new_image[i][j] = image[i-1][j-1]
		}
	}
	return new_image
}
func crop(image_agrandit [][]int, x int, y int, size int) [][]int { // récupére un carré de l'image originale centré en x,y et de dimension size*size
	img_result := slice2D(size, size)

	for ligne := 0; ligne < size; ligne++ { //On remplie le carré par les valleurs correspondantes
		for collone := 0; collone < size; collone++ {
			img_result[ligne][collone] = image_agrandit[y+ligne][x+collone]
		}
	}
	return img_result
}

func sum2D(kernel [][]int) int {
	result := 0
	for i := range kernel {
		for j := range kernel[0] {
			result += kernel[i][j]
		}
	}
	return result
}

func computeconvolution(result_array [][]int, image_agrandit [][]int, kernel [][]int, x int, y int) {
	size := len(kernel)
	croped_image := crop(image_agrandit, x, y, size)
	result := 0
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			result += kernel[i][j] * croped_image[size-i-1][size-j-1] // On effectue le calcul de la convolution on ajoutant les elements opposé entre le filtre et l'image
		}
	}
	somme := sum2D(kernel) //On normalise le resultat par la somme des coefficients du filtre si le filtre le permet
	if somme != 0 {
		result_array[y][x] = result / somme
	} else {
		result_array[y][x] = result
	}
}

func line_compute(result_array [][]int, image_agrandit [][]int, kernel [][]int, y int, lenx int, leny int, nb_ligne int) {
	for i := y; i < y+nb_ligne || i < leny; i++ {
		for j := 0; j < lenx; j++ {
			computeconvolution(result_array, image_agrandit, kernel, j, i)
		}
	}
}

func convolute(image_array [][]int, kernel [][]int) [][]int { //Fonction a appeler pour effectuer la convolution d'une image et d'un filtre
	leny := len(image_array)
	lenx := len(image_array[0])
	image_agrandit := agrandit(image_array) //On traite l'image pour rajouter des 0 sur les bordures
	result := slice2D(leny, lenx)
	nb_ligne := 10 //On définit le nombre de ligne que va calculer chaque go routine

	//maxGoroutines := 10000
	//guard := make(chan int, maxGoroutines) // On crée un channel qui va limiter le nombre de go routine qui tourneront en même temps

	for i := 0; i < leny; i += nb_ligne {
		//guard <- 1 //On rempli le channel lors du lancement de la routine
		//go func(j int, i int) {
		go line_compute(result, image_agrandit, kernel, i, lenx, leny, nb_ligne)
		//<-guard //On vide le channel à la fin de la tache
		//}(j, i)
	}
	return result
}

func main() {
	start := time.Now()

	image := slice2D(1920, 1080) //Creation du 2D image pour tester
	for i := 0; i < len(image); i++ {
		for j := 0; j < len(image[0]); j++ {
			image[i][j] = j + 10*i
		}
	}
	kernel := slice2D(3, 3) //Creation du kernel pour tester
	kernel[1][1] = 1

	final := convolute(image, kernel)
	elapsed := time.Since(start)
	fmt.Printf("Temps: %s", elapsed)
	final = final
}
