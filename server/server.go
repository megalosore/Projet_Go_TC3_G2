package main

import (
	"bufio"
	"encoding/gob"
	"flag"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var routineNb int
var inputChannel chan *toCompute

func getArgs() int {
	usageString := "Usage: go run server.go [-C=NumberRoutine] <portnumber>\n"

	flagNbroutine := flag.Int("C", runtime.NumCPU(), "Number of go routine per client")
	flag.Parse()
	routineNb = *flagNbroutine

	if len(flag.Args()) != 1 {
		fmt.Printf(usageString)
		os.Exit(1)

	} else {
		fmt.Printf("#DEBUG Arg Port Number : %s\n", flag.Arg(0))
		portNumber, err := strconv.Atoi(flag.Arg(0))
		if err != nil {
			fmt.Printf(usageString)
			os.Exit(1)
		} else {
			return portNumber
		}
	}
	return -1
}

func handleConnection(connection net.Conn, connum int) {
	defer connection.Close()
	connReader := bufio.NewReader(connection)

	for {
		inputLine, err := connReader.ReadString('\n')
		if err != nil {
			fmt.Printf("Connection %d encountered an error.\n", connum)
			fmt.Printf("Error: %s\n", err.Error())
			break
		}

		argsString := strings.TrimSuffix(inputLine, "\n")
		argsList := strings.Split(argsString, "\\")
		url := argsList[0]
		kernelType := argsList[1]
		thresholdValue := argsList[2]

		fmt.Printf("#DEBUG Connection %d Args: url=%s, alg=%s threshold=%s\n", connum, url, kernelType, thresholdValue)

		img, err := loadImgFromURL(url)
		if err != nil {
			fmt.Printf("Connection %d encountered an error.\n", connum)
			fmt.Printf("Error: %s\n", err.Error())
			break
		}

		inputSlice := imgToSlice(img)
		outputSlice := slice2D(len(inputSlice), len(inputSlice[0]))
		var threshold float64
		var kernel1 [][]int16
		var kernel2 [][]int16
		doubleKernel := false
		outputChannel := make(chan bool)

		// choix du/des kernels à utiliser
		switch kernelType {
		case "sobel":
			doubleKernel = true
			kernel1 = [][]int16{
				{-1, 0, 1},
				{-2, 0, 2},
				{-1, 0, 1},
			}
			kernel2 = [][]int16{
				{-1, -2, -1},
				{0, 0, 0},
				{1, 2, 1},
			}
			if thresholdValue != "" {
				threshold, err = strconv.ParseFloat(thresholdValue, 8)
				if err != nil {
					threshold = 0.1 // Default threshold value for Sobel
				}
			} else {
				threshold = 0.1 // Default threshold value for Sobel
			}

		case "prewit":
			doubleKernel = true
			kernel1 = [][]int16{
				{-1, 0, 1},
				{-1, 0, 1},
				{-1, 0, 1},
			}
			kernel2 = [][]int16{
				{-1, -1, -1},
				{0, 0, 0},
				{1, 1, 1},
			}
			if thresholdValue != "" {
				threshold, err = strconv.ParseFloat(thresholdValue, 8)
				if err != nil {
					threshold = 0.07 // Default threshold value for Prewit
				}
			} else {
				threshold = 0.07 // Default threshold value for Prewit
			}

		default:
			// On utilise le Laplacien par défaut sinon
			kernel1 = [][]int16{
				{0, -1, 0},
				{-1, 4, -1},
				{0, -1, 0},
			}
			kernel2 = nil
			if thresholdValue != "" {
				threshold, err = strconv.ParseFloat(thresholdValue, 8)
				if err != nil {
					threshold = 0.3 // Default threshold value for Laplacian
				}
			} else {
				threshold = 0.3 // Default threshold value for Laplacian
			}
		}
		start := time.Now()
		feedInput(inputSlice, outputSlice, doubleKernel, kernel1, kernel2, threshold, outputChannel)
		nbReceived := 0
		for nbReceived < routineNb {
			_ = <-outputChannel
			nbReceived++
		}
		elapsed := time.Since(start)
		fmt.Printf("Computation time: %s\n", elapsed)

		outputImg := sliceToImg(outputSlice)
		encoder := gob.NewEncoder(connection)
		encoder.Encode(outputImg)
		break
	}
}

func loadImgFromURL(url string) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func main() {
	port := getArgs()
	fmt.Printf("#DEBUG Creating TCP Server on port %d\n", port)
	portString := fmt.Sprintf(":%s", strconv.Itoa(port))
	fmt.Printf("#DEBUG Number of go routines: %d\n", routineNb)

	ln, err := net.Listen("tcp", portString)
	if err != nil {
		fmt.Printf("Error: Could not create listener.\n")
		panic(err)
	}

	inputChannel = make(chan *toCompute)
	launchWorkers()
	connum := 1

	for {
		fmt.Printf("#DEBUG Accepting next connection\n")
		conn, errconn := ln.Accept()

		if errconn != nil {
			fmt.Printf("Error when accepting next connection.\n")
			panic(errconn)
		}

		go handleConnection(conn, connum)
		connum += 1
	}
}
