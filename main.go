package main

import (
	"flag"
	"fmt"
	"log"
	"sort"

	"fake.com/PortScanner/port"
)

const (
	DEFAULT_PROTOCOL = "TCP"
	DEFAULT_HOSTNAME = "127.0.0.1"
	MIN_START_PORT   = 1
	MAX_END_PORT     = 65535
)

func main() {
	protocol := flag.String("protocol", DEFAULT_PROTOCOL, "Either TCP or UDP")
	hostname := flag.String("host", DEFAULT_HOSTNAME, "Host IP address")
	startPort := flag.Int("start", MIN_START_PORT, "Port to start when scanning")
	endPort := flag.Int("end", MAX_END_PORT, "Port to end when scanning")
	flag.Parse()

	var openPorts []int

	if *startPort < MIN_START_PORT {
		err := fmt.Errorf("Invalid start port: port cannot be < %d\n", MIN_START_PORT)
		log.Fatal(err)
	} else if *endPort > MAX_END_PORT {
		err := fmt.Errorf("Invalid end port: port cannot be > %d\n", MAX_END_PORT)
		log.Fatal(err)
	}

	if *protocol == "UDP" {
		openPorts = port.RunWideUDPScan(*startPort, *endPort, *hostname)
	} else if *protocol == "TCP" {
		openPorts = port.RunWideTCPScan(*startPort, *endPort, *hostname)
	} else {
		err := fmt.Errorf("Invalid Protocol: %s is not valid protocol (only TCP and UDP are valid)\n", *protocol)
		log.Fatal(err)
	}

	sort.Ints(openPorts)
	printPorts(*protocol, *hostname, openPorts)
}

func printPorts(protocol string, hostname string, openPorts []int) {
	if len(openPorts) > 0 {
		fmt.Printf("-------------- Open %s Ports on Host %q --------------\n", protocol, hostname)
		for _, port := range openPorts {
			fmt.Println(port)
		}
		fmt.Printf("-----------------------------------------------------------------\n")
	} else {
		fmt.Printf("-------------- No Open %s Port on Host %q in given range --------------\n", protocol, hostname)
	}
}
