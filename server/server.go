package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"image"
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

		url := strings.TrimSuffix(inputLine, "\n")
		fmt.Printf("#DEBUG %d RCV |%s|\n", connum, url)

		start := time.Now()
		img, err := loadImgFromURL(url)
		if err != nil {
			fmt.Printf("#DEBUG %d RCV ERROR no panic, just a client\n", connum)
			fmt.Printf("Error :|%s|\n", err.Error())
			break
		}

		imgConverted := imgToSlice(img)
		kernel := [][]int16{
			{0, -1, 0},
			{-1, 4, -1},
			{0, -1, 0},
		}

		final := convolute(imgConverted, kernel)
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
