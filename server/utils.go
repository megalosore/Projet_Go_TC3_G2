package main

import (
	"fmt"
	"image"
	"image/color"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
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
func crop(inputSlice [][]int16, x int, y int, size int) [][]int16 {
	outputSlice := slice2D(size, size)

	// On remplit le carré par les valeurs correspondantes
	for ligne := 0; ligne < size; ligne++ {
		for colonne := 0; colonne < size; colonne++ {
			outputSlice[ligne][colonne] = inputSlice[y-1+ligne][x-1+colonne]
		}
	}
	return outputSlice
}

// Crée une version entourée de 0 de l'image originale pour traiter les cas des x,y en bordure
func fillBorders(slice [][]int16) [][]int16 {
	newImage := slice2D(len(slice)+2, len(slice[0])+2)
	for i := 1; i < len(newImage)-1; i++ {
		for j := 1; j < len(newImage[0])-1; j++ {
			newImage[i][j] = slice[i-1][j-1]
		}
	}
	return newImage
}

// Crée un double slice de dimension précisée (y=ligne, x=colonne) rempli de 0
func slice2D(lenY int, lenX int) [][]int16 {
	doubleSlice := make([][]int16, lenY)
	for i := range doubleSlice {
		doubleSlice[i] = make([]int16, lenX)
	}
	return doubleSlice
}

// Convertit un slice 2D en une image en nuances de gris
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

// Trace un graphe et l'enregistre dans un répertoire spécifique
func traceGraph(graphName string, xData []int, yData []float64, xAxisName string, yAxisName string, baseName string) {
	items := make([]opts.LineData, 0)
	for i := 0; i < len(yData); i++ {
		items = append(items, opts.LineData{Value: yData[i], Symbol: "circle", SymbolSize: 5})
	}
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: graphName}),
		charts.WithYAxisOpts(opts.YAxis{Name: yAxisName, Scale: true}),
		charts.WithXAxisOpts(opts.XAxis{Name: xAxisName, Scale: true}),
		charts.WithDataZoomOpts(opts.DataZoom{Type: "slider", Start: float32(xData[0]), End: float32(xData[len(xData)-1]), XAxisIndex: []int{0}}),
	)
	line.SetXAxis(xData).AddSeries("", items)

	page := components.NewPage()
	page.AddCharts(line)

	folderName := baseName + "_results"
	// On crée un répertoire pour enregistrer les graphes s'il n'existe pas déjà
	if _, err := os.Stat(folderName + "/"); os.IsNotExist(err) {
		err = os.Mkdir(folderName, 0755)
		if err != nil {
			fmt.Printf("Error while creating %s: check the permissions of the current directory.\n", folderName)
		}
	}
	// On enregistre le graphe sans écraser les résultats précédents
	files, _ := ioutil.ReadDir(folderName + "/")
	filePath := folderName + "/" + baseName + "_" + strconv.Itoa(len(files)+1) + ".html"
	f, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Error while saving graph.\n")
	} else {
		fmt.Printf("Graph result successfully saved at \"%s\".\n", filePath)
	}
	page.Render(io.MultiWriter(f))
}

// Mesure la complexité en temps de la convolution
func complexityBenchmark() {
	fmt.Printf("Starting benchmark\n")
	var imgSizes []int
	// On détermine ici le nombre et l'abscisse des points
	for i := 10; i <= 10000; i += 10 {
		imgSizes = append(imgSizes, i)
	}
	times := make([]float64, len(imgSizes))
	outputChannel := make(chan bool)
	launchWorkers()

	for j := range imgSizes {
		fmt.Printf("Computation %d/%d\n", j+1, len(imgSizes))
		start := time.Now()
		inputs := generateInputs(slice2D(imgSizes[j], imgSizes[j]), slice2D(imgSizes[j], imgSizes[j]), true, slice2D(3, 3), slice2D(3, 3), 0.5, outputChannel)
		go feedInput(inputs)
		nbReceived := 0
		for nbReceived < len(inputs) {
			_ = <-outputChannel
			nbReceived++
		}
		elapsed := time.Since(start)
		times[j] = float64(elapsed.Milliseconds())
	}
	fmt.Printf("Benchmark finished.\n")
	traceGraph("Complexité en temps", imgSizes, times, "Taille image", "Temps (ms)", "complexity_benchmark")
}
