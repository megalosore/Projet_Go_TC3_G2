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

func getArgs() int {

	if len(os.Args) != 2 {
		fmt.Printf("Usage: 	go run client.go <portnumber>\n")
		os.Exit(1)
	} else {
		fmt.Printf("#DEBUG ARGS Port Number : %s\n", os.Args[1])
		portNumber, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Printf("Usage: go run client.go <portnumber>\n")
			os.Exit(1)
		} else {
			return portNumber
		}

	}
	//Should never be reached
	return -1
}

func writeImg(img *image.Gray, imgName string) error {
	f, err := os.Create(imgName + ".png")
	if err != nil {
		return err
	}
	defer f.Close()

	err = png.Encode(f, img)
	return err
}

func main() {
	port := getArgs()
	fmt.Printf("#DEBUG DIALING TCP Server on port %d\n", port)
	portString := fmt.Sprintf("127.0.0.1:%s", strconv.Itoa(port))
	fmt.Printf("#DEBUG MAIN PORT STRING |%s|\n", portString)

	conn, err := net.Dial("tcp", portString)
	if err != nil {
		fmt.Printf("#DEBUG MAIN could not connect\n")
		os.Exit(1)
	} else {
		url := "https://f.hellowork.com/blogdumoderateur/2013/02/nyan-cat-gif-1.gif"
		io.WriteString(conn, url+"\n")
		defer conn.Close()
		decoder := gob.NewDecoder(conn)
		var img image.Gray
		decoder.Decode(&img)
		writeImg(&img, "output")
	}

}
