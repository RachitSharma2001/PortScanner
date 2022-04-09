package main

import (
	"fmt"
	"net"
)

func main() {
	port := 120
	server, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Printf("Error listening at port %d: %v\n", port, err)
	} else {
		for {
			_, err := server.Accept()
			if err != nil {
				fmt.Printf("Error while connecting: %v\n", err)
			}
		}
	}
}
