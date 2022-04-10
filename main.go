package main

import (
	port "fake.com/PortScanner/port"
)

func main() {
	port.RunWideUDPScan(1, 65535)
	port.RunWideTCPScan(1, 65535)
}
