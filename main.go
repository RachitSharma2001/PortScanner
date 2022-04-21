package main

import (
	"fmt"
	"sort"

	port "fake.com/PortScanner/port"
)

const (
	DEFAULT_PROTOCOL = "TCP"
	DEFAULT_HOSTNAME = "127.0.0.1"
)

func main() {
	protocol := DEFAULT_PROTOCOL
	hostname := DEFAULT_HOSTNAME
	openTcpPorts := port.RunWideTCPScan(1, 65535, hostname)
	sort.Ints(openTcpPorts)
	printPorts(protocol, hostname, openTcpPorts)
}

func printPorts(protocol string, hostname string, openPorts []int) {
	fmt.Printf("-------------- Open %s Ports on Host %q --------------\n", protocol, hostname)
	for _, port := range openPorts {
		fmt.Println(port)
	}
	fmt.Printf("-----------------------------------------------------------------\n")
}
