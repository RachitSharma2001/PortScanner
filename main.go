package main

import (
	"fmt"

	port "fake.com/PortScanner/port"
)

func main() {
	/*openPorts := port.RunWideUDPScan(1, 65535, "0.0.0.0")
	fmt.Println("Open UDP ports:")
	fmt.Println(openPorts)*/

	openPorts := port.RunWideTCPScan(1, 65535, "0.0.0.0")
	fmt.Println("Open TCP ports:")
	fmt.Println(openPorts)
	//port.RunWideTCPScan(1, 65535, "localhost")
}
