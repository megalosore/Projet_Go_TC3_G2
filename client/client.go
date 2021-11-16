package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"image"
	"image/png"
	"io"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func getArgs() (int, string, string, string, string) {
	destPtr := flag.String("D", "", "Destination Path for the output file") // Mise en place des flags pour préciser les arguments optionnels
	algPtr := flag.String("A", "", "Name of the algorithme that will be used by the server")
	thresholdPtr := flag.String("S", "", "Value of the threshold that will be used by the server")
	flag.Parse()

	destinationPath := *destPtr // Attribution des valeurs des arguments optionnels
	alg := *algPtr
	threshold := *thresholdPtr
	imageURL := flag.Arg(1)

	if len(flag.Args()) != 2 {
		fmt.Printf("Usage: go run client.go [-D=destination_path] [-A=algorithm] [-S=threshold_value] <server_portnumber> <image_url>\n")
		os.Exit(2)
	}

	fmt.Printf("#DEBUG ARGS portNumber : %s\n", flag.Arg(0))
	portNumber, errPort := strconv.Atoi(flag.Arg(0))
	if errPort != nil {
		fmt.Printf("Error: incorrect port number\n")
		fmt.Printf("Usage: go run client.go [-D=destination_path] [-A=algorithm] [-S=threshold_value] <server_portnumber> <image_url>\n")
		os.Exit(2)
	}

	fmt.Printf("#DEBUG ARGS URL : %s\n", flag.Arg(1))
	_, errUrl := url.ParseRequestURI(imageURL) // check if the URL respect HTTP URL format
	if errUrl != nil {
		fmt.Printf("Error: invalid URL\n")
		fmt.Printf("Usage: go run client.go [-D=destination_path] [-A=algorithm] [-S=threshold_value] <server_portnumber> <image_url>\n")
		os.Exit(2)
	}
	if threshold != "" {
		thresholdFloat, errThreshold := strconv.ParseFloat(threshold, 8)
		if errThreshold != nil || thresholdFloat < 0 || thresholdFloat > 1 { //Check if the threshold value is a number between 0 and 1
			fmt.Printf("Error: incorrect threshold value\n")
			fmt.Printf("Please enter a threshold value between 0 and 1\n")
			fmt.Printf("Usage: go run client.go [-D=destination_path] [-A=algorithm] [-S=threshold_value] <server_portnumber> <image_url>\n")
			os.Exit(2)
		}
	}
	fmt.Printf("#DEBUG ARGS DestinationPath : %s\n", destinationPath)
	fmt.Printf("#DEBUG ARGS Algorithme : %s\n", alg)
	fmt.Printf("#DEBUG ARGS thresholdValue : %s\n", threshold)
	return portNumber, imageURL, destinationPath, alg, threshold
}

func writeImg(img *image.Gray, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	err = png.Encode(f, img)
	return err
}

func generateImgPath(actualPath string) string {
	if actualPath == "" {
		return "output.png"
	}

	if strings.HasSuffix(actualPath, "/") || strings.HasSuffix(actualPath, "\\") {
		return actualPath + "output.png"
	}

	if len(actualPath) > 4 {
		suffix := actualPath[len(actualPath)-4:]
		suffix = strings.ToLower(suffix)
		switch suffix {
		case ".png", ".jpg", ".gif":
			return actualPath[:len(actualPath)-4] + ".png"
		}
	}

	return actualPath + ".png"
}

func main() {
	portNumber, imageURL, destinationPath, alg, threshold := getArgs()
	destinationPath = generateImgPath(destinationPath)
	sendValue := imageURL + "\\" + alg + "\\" + threshold + "\n"

	fmt.Printf("#DEBUG DIALING TCP Server on port %d\n", portNumber)
	portString := fmt.Sprintf("127.0.0.1:%s", strconv.Itoa(portNumber))
	fmt.Printf("#DEBUG MAIN PORT STRING |%s|\n", portString)
	conn, err := net.Dial("tcp", portString)

	if err != nil {
		fmt.Printf("#DEBUG MAIN could not connect\n")
		os.Exit(3)
	} else {
		io.WriteString(conn, sendValue) // Send data to the server
		defer conn.Close()
		decoder := gob.NewDecoder(conn)
		var img image.Gray
		err := decoder.Decode(&img)
		if err != nil {
			panic(err)
		}
		err = writeImg(&img, destinationPath)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Image successfully saved at %s\n", destinationPath)
	}
}
