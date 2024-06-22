package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println("Could not listen to port: ", err)
		return
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Could not listen to request: ", err)
			return
		}
		defer conn.Close()

		buf := make([]byte, 1024)

		_, err = conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error reading from client", err.Error())
			os.Exit(1)
		}

		conn.Write([]byte("+OK\r\n"))
	}
}
