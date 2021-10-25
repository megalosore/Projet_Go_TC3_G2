package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func getArgs() int {

	if len(os.Args) != 2 {
		fmt.Printf("Usage: go run server.go <portnumber>\n")
		os.Exit(1)
	} else {
		fmt.Printf("#DEBUG ARGS Port Number : %s\n", os.Args[1])
		portNumber, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Printf("Usage: go run server.go <portnumber>\n")
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
			fmt.Printf("#DEBUG %d RCV ERROR no panic, just a client\n", connum)
			fmt.Printf("Error :|%s|\n", err.Error())
			break
		}

		argsString := strings.TrimSuffix(inputLine, "\n")
		argsList := strings.Split(argsString, "\\")
		url := argsList[0]
		kernelType := argsList[1]
		seuilValue := argsList[2]

		fmt.Printf("#DEBUG %d RCV |%s| |%s| |%s|\n", connum, url, kernelType, seuilValue)

		start := time.Now()
		img, err := loadImgFromURL(url)
		if err != nil {
			fmt.Printf("#DEBUG %d RCV ERROR no panic, just a client\n", connum)
			fmt.Printf("Error :|%s|\n", err.Error())
			break
		}

		var final [][]int16 //On initialise la valeur qui reçoit le resultat de nos calculs
		var seuil float64   //On initialise la valeur qui recevra le seuil précisé ou pas par le client

		switch kernelType {
		case "sobel": //Si le client spécifie le filtre de sobel on l'utilise
			kernel1 := [][]int16{
				{-1, 0, 1},
				{-2, 0, 2},
				{-1, 0, 1},
			}
			kernel2 := [][]int16{
				{-1, -2, -1},
				{0, 0, 0},
				{1, 2, 1},
			}
			if seuilValue != "" {
				seuil, err = strconv.ParseFloat(seuilValue, 8)
				if err != nil {
					seuil = 0.1 //Default seuil value for Sobel
				}
			} else {
				seuil = 0.1 //Default seuil value for Sobel
			}
			imgConverted := imgToSlice(img)
			final = convoluteDouble(imgConverted, kernel1, kernel2, seuil)

		case "prewit": //Si le client spécifie le filtre de prewit on l'utilise
			kernel1 := [][]int16{
				{-1, 0, 1},
				{-1, 0, 1},
				{-1, 0, 1},
			}
			kernel2 := [][]int16{
				{-1, -1, -1},
				{0, 0, 0},
				{1, 1, 1},
			}
			if seuilValue != "" {
				seuil, err = strconv.ParseFloat(seuilValue, 8)
				if err != nil {
					seuil = 0.07 //Default seuil value for Prewit
				}
			} else {
				seuil = 0.07 //Default seuil value for Prewit
			}
			imgConverted := imgToSlice(img)
			final = convoluteDouble(imgConverted, kernel1, kernel2, seuil)

		default:
			//On utilise le Laplacien par défaut sinon
			kernel := [][]int16{
				{0, -1, 0},
				{-1, 4, -1},
				{0, -1, 0},
			}
			if seuilValue != "" {
				seuil, err = strconv.ParseFloat(seuilValue, 8)
				if err != nil {
					seuil = 0.3 //Default seuil value for Laplacian
				}
			} else {
				seuil = 0.3 //Default seul value for Laplacian
			}
			imgConverted := imgToSlice(img)
			final = convolute(imgConverted, kernel, seuil)
		}

		finalImage := sliceToImg(final)
		elapsed := time.Since(start)
		fmt.Printf("Temps : %s\n", elapsed)
		encoder := gob.NewEncoder(connection)
		encoder.Encode(finalImage)
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
	fmt.Printf("#DEBUG MAIN Creating TCP Server on port %d\n", port)
	portString := fmt.Sprintf(":%s", strconv.Itoa(port))
	fmt.Printf("#DEBUG MAIN PORT STRING |%s|\n", portString)

	ln, err := net.Listen("tcp", portString)
	if err != nil {
		fmt.Printf("#DEBUG MAIN Could not create listener\n")
		panic(err)
	}

	connum := 1

	for {
		fmt.Printf("#DEBUG MAIN Accepting next connection\n")
		conn, errconn := ln.Accept()

		if errconn != nil {
			fmt.Printf("DEBUG MAIN Error when accepting next connection\n")
			panic(errconn)
		}

		go handleConnection(conn, connum)
		connum += 1
	}
}
