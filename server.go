package main

import (
	"fmt"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("error!")
	}

	for {
		_, err := ln.Accept()
		if err != nil {
			fmt.Println("error!")
		} else {
			fmt.Println("Got a connection!")
		}
	}
}
