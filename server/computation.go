package main

import (
	"math"
)

type toCompute struct {
	enlargedSlice [][]int16
	outputSlice   [][]int16
	lenX          int
	doubleKernel  bool
	kernel1       [][]int16
	kernel2       [][]int16
	threshold     float64
	startingLine  int
	lineNumber    int
	outputChannel chan bool
	killSignal    bool
}

// Effectue la convolution pour 1 pixel avec un seul kernel
func computeConvolutionSimple(input *toCompute, x int, y int) {
	size := len(input.kernel1)
	// On récupère une version size*size entourant le pixel que l'on veut traiter
	croppedSlice := crop(input.enlargedSlice, x+1, y+1, size)
	result := int16(0)
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			// On effectue le calcul de la convolution on ajoutant les éléments opposés entre le filtre et l'image cropée
			result += input.kernel1[i][j] * croppedSlice[size-i-1][size-j-1]
		}
	}
	if result < 0 { //On borne notre resultat entre 0 et 255
		result = -1 * result
	}
	if result > 255 {
		result = 255
	}

	// On normalise le résultat par la somme des coefficients du filtre si le filtre le permet
	sum := sum2D(input.kernel1)
	if sum != 0 {
		input.outputSlice[y][x] = result / sum
	} else {
		//seuillage pour les detections de contours
		if result < int16(255*input.threshold) {
			input.outputSlice[y][x] = 0
		} else {
			input.outputSlice[y][x] = 255
		}
	}
}

// Effectue la convolution pour 1 pixel avec deux kernels
func computeConvolutionDouble(input *toCompute, x int, y int) {
	size := len(input.kernel1)
	// On récupère une version size*size entourant le pixel que l'on veut traiter
	croppedSlice := crop(input.enlargedSlice, x+1, y+1, size)
	result1 := int16(0)
	result2 := int16(0)
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			// On effectue le calcul de la convolution on ajoutant les éléments opposés entre le filtre et l'image cropée
			// On implémente la convolution de Sobel
			result1 += input.kernel1[i][j] * croppedSlice[size-i-1][size-j-1]
			result2 += input.kernel2[i][j] * croppedSlice[size-i-1][size-j-1]
		}
	}
	result := math.Sqrt(math.Pow(float64(result1), 2)+math.Pow(float64(result2), 2)) / 4
	if result < input.threshold*255 { //Seuillage
		input.outputSlice[y][x] = 0
	} else {
		input.outputSlice[y][x] = 255
	}
}

// Calcule la convolution d'une partie de l'image et du/des kernels
func worker() {
	for {
		input := <-inputChannel

		if input.killSignal {
			input.outputChannel <- true
			break
		}

		for i := input.startingLine; i < input.startingLine+input.lineNumber; i++ {
			for j := 0; j < input.lenX; j++ {
				if input.doubleKernel {
					computeConvolutionDouble(input, j, i)
				} else {
					computeConvolutionSimple(input, j, i)
				}
			}
		}
		input.outputChannel <- true
	}
}

func launchWorkers() {
	for i := 0; i < routineNb; i++ {
		go worker()
	}
}

// Envoie un signal d'arrêt aux goroutines et s'assure qu'elles sont bien tuées
func killWorkers(nbToKill int) {
	outputChannel := make(chan bool)
	for i := 0; i < nbToKill; i++ {
		inputChannel <- &toCompute{killSignal: true, outputChannel: outputChannel}
	}
	nbReceived := 0
	for nbReceived < nbToKill {
		_ = <-outputChannel
		nbReceived++
	}
}

// Découpe le slice en blocks à distribuer aux workers
func generateInputs(inputSlice [][]int16, outputSlice [][]int16, doubleKernel bool, kernel1 [][]int16, kernel2 [][]int16, threshold float64, outputChannel chan bool) []*toCompute {
	lenY := len(inputSlice)
	lenX := len(inputSlice[0])

	// On traite l'image pour rajouter des 0 sur les bordures
	enlargedSlice := fillBorders(inputSlice)

	// On découpe l'image en chunks pour correspondre au nombre de goroutines
	q := float64(lenY) / float64(routineNb)
	toAdd := 0
	lineNb := 1
	if q > 1 {
		lineNb = int(q)
		toAdd = int(math.Round((q - float64(lineNb)) * float64(routineNb)))
	}

	var inputs []*toCompute
	var n int
	chunkNumber := 0
	for i := 0; i < lenY; i += lineNb {
		n = lineNb
		if toAdd != 0 {
			n++
			i++
			toAdd--
		}
		inputs = append(inputs, &toCompute{enlargedSlice, outputSlice, lenX, doubleKernel,
			kernel1, kernel2, threshold, i, n, outputChannel, false})
		chunkNumber++
	}
	return inputs
}

// Remplit le channel d'input
func feedInput(inputs []*toCompute) {
	for _, input := range inputs {
		inputChannel <- input
	}
}
