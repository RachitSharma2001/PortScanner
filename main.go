package main

import (
	"flag"
	"fmt"
	"log"
	"sort"
	"strings"

	"fake.com/PortScanner/port"
)

const (
	DEFAULT_PROTOCOL = "TCP"
	DEFAULT_HOSTNAME = "127.0.0.1"
	MIN_START_PORT   = 1
	MAX_END_PORT     = 65535
)

var (
	protocol  *string
	hostname  *string
	startPort *int
	endPort   *int
)

func main() {
	parseUserEnteredFlags()
	validatePortRange()
	validateProtocol()
	openPorts := scanForOpenPorts()
	sort.Ints(openPorts)
	printPorts(openPorts)
}

func parseUserEnteredFlags() {
	protocol = flag.String("protocol", DEFAULT_PROTOCOL, "Either TCP or UDP")
	hostname = flag.String("host", DEFAULT_HOSTNAME, "Host IP address")
	startPort = flag.Int("start", MIN_START_PORT, "Port to start when scanning")
	endPort = flag.Int("end", MAX_END_PORT, "Port to end when scanning")
	flag.Parse()
}

func validatePortRange() {
	if *startPort < MIN_START_PORT {
		err := fmt.Errorf("Invalid start port: port cannot be < %d\n", MIN_START_PORT)
		log.Fatal(err)
	} else if *endPort > MAX_END_PORT {
		err := fmt.Errorf("Invalid end port: port cannot be > %d\n", MAX_END_PORT)
		log.Fatal(err)
	} else if *startPort > *endPort {
		err := fmt.Errorf("Invalid port range: start port must be <= end port\n")
		log.Fatal(err)
	}
}

func validateProtocol() {
	if *protocol != "UDP" && *protocol != "TCP" {
		err := fmt.Errorf("Invalid Protocol: %s is not valid protocol (only TCP and UDP are valid)\n", *protocol)
		log.Fatal(err)
	}
}

func scanForOpenPorts() []int {
	var openPorts []int
	if *protocol == "UDP" {
		openPorts = port.RunWideUDPScan(*startPort, *endPort, *hostname)
	} else {
		openPorts = port.RunWideTCPScan(*startPort, *endPort, *hostname)
	}
	return openPorts
}

func printPorts(openPorts []int) {
	if len(openPorts) > 0 {
		outputAllOpenPorts(openPorts)
	} else {
		outputPortIsEmpty()
	}
}

func outputAllOpenPorts(openPorts []int) {
	openPortMsg := fmt.Sprintf("----------- Open %s Ports on Host %q from %d - %d -----------\n", *protocol, *hostname, *startPort, *endPort)
	fmt.Printf(openPortMsg)
	for _, port := range openPorts {
		fmt.Println(port)
	}
	fmt.Println(strings.Repeat("-", len(openPortMsg)-1))
}

func outputPortIsEmpty() {
	fmt.Printf("-------------- No Open %s Port on Host %q from %d - %d --------------\n", *protocol, *hostname, *startPort, *endPort)
}
