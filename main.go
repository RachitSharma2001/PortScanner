package main

import (
	"fmt"

	port "fake.com/PortScanner/port"
)

func main() {
	/*openUdpPorts := port.RunWideUDPScan(1, 65535, "127.0.0.1")
	printPorts("UDP", openUdpPorts)*/

	openTcpPorts := port.RunWideTCPScan(1, 65535, "127.0.0.1")
	printPorts("TCP", openTcpPorts)
}

func printPorts(protocol string, openPorts []int) {
	fmt.Printf("Open %s ports: %v \n", protocol, openPorts)
}
