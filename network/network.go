package network

import "net"

func test() {
	addy, err := net.ResolveUDPAddr("udp", "127.0.0.1")
}
