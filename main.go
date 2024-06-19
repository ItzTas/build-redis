package main

import (
	"fmt"
	"net"
)

func main() {
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println("Could not listen to port: ", err)
		return
	}
	defer l.Close()

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Could not listen to request: ", err)
		return
	}
	defer conn.Close()
	for {
		buf := make([]byte, 1024)

		_, err := conn.Read(buf)
		if err != nil {
			fmt.Println("")
			return
		}
	}
}
