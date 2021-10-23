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
)

func getArgs() (int, string, string) {

	if len(os.Args) != 4 {
		fmt.Printf("Usage: 	go run client.go <serveur_portnumber> <URL_IMG> <destinationPath>\n")
		os.Exit(1)
	} else {
		fmt.Printf("#DEBUG ARGS Port Number : %s\n", os.Args[1])
		portNumber, err_port := strconv.Atoi(os.Args[1])
		fmt.Printf("#DEBUG ARGS URL_IMG : %s\n", os.Args[2])
		URL_IMG := os.Args[2]
		fmt.Printf("#DEBUG ARGS destinationPath : %s\n", os.Args[3])
		destinationPath := os.Args[3]
		if err_port != nil {
			fmt.Printf("Usage: go run client.go <serveur_portnumber> <URL_IMG> <destinationPath>\n")
			os.Exit(1)
		} else {
			return portNumber, URL_IMG, destinationPath
		}

	}
	//Should never be reached
	return -1, "", ""
}

func writeImg(img *image.Gray, imgName string) error {
	f, err := os.Create(imgName)
	if err != nil {
		return err
	}
	defer f.Close()

	err = png.Encode(f, img)
	return err
}

func main() {
	port, url, destinationPath := getArgs()
	fmt.Printf("#DEBUG DIALING TCP Server on port %d\n", port)
	portString := fmt.Sprintf("127.0.0.1:%s", strconv.Itoa(port))
	fmt.Printf("#DEBUG MAIN PORT STRING |%s|\n", portString)

	conn, err := net.Dial("tcp", portString)
	if err != nil {
		fmt.Printf("#DEBUG MAIN could not connect\n")
		os.Exit(1)
	} else {
		io.WriteString(conn, url+"\n")
		defer conn.Close()
		decoder := gob.NewDecoder(conn)
		var img image.Gray
		decoder.Decode(&img)
		writeImg(&img, destinationPath)
	}

}
