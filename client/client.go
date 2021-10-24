package main

import (
	"encoding/gob"
	"fmt"
	"image"
	"image/png"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

func getArgs() (int, string, string) {
	portNumber := 0
	imageURL := ""
	destinationPath := ""

	if len(os.Args) < 3 || len(os.Args) > 4 {
		fmt.Printf("Usage: go run client.go <server_portnumber> <image_url> [destinationPath]\n")
		os.Exit(2)
	} else {
		fmt.Printf("#DEBUG ARGS portNumber : %s\n", os.Args[1])
		var errPort error
		portNumber, errPort = strconv.Atoi(os.Args[1])
		if errPort != nil {
			fmt.Printf("Error: incorrect port number")
			fmt.Printf("Usage: go run client.go <server_portnumber> <image_url> [destinationPath]\n")
			os.Exit(2)
		}

		fmt.Printf("#DEBUG ARGS imageURL : %s\n", os.Args[2])
		imageURL = os.Args[2]

		if len(os.Args) == 4 {
			fmt.Printf("#DEBUG ARGS destinationPath : %s\n", os.Args[3])
			destinationPath = os.Args[3]
		}
	}
	return portNumber, imageURL, destinationPath
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
	port, url, destinationPath := getArgs()
	destinationPath = generateImgPath(destinationPath)

	fmt.Printf("#DEBUG DIALING TCP Server on port %d\n", port)
	portString := fmt.Sprintf("127.0.0.1:%s", strconv.Itoa(port))
	fmt.Printf("#DEBUG MAIN PORT STRING |%s|\n", portString)
	conn, err := net.Dial("tcp", portString)

	if err != nil {
		fmt.Printf("#DEBUG MAIN could not connect\n")
		os.Exit(3)
	} else {
		io.WriteString(conn, url+"\n")
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
		fmt.Printf("Image successfully saved at %s", destinationPath)
	}
}
