package main

import (
	"fmt"
	"net"
)

func main() {
	fmt.Println("Listening on port :6379")

	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		defer conn.Close()
		resp := NewResp(conn)
		_, err = resp.Read()
		if err != nil {
			fmt.Println(err)
			return
		}

		writer := NewWriter(conn)
		writer.Write(Value{typ: "string", str: "OK"})
	}
}
