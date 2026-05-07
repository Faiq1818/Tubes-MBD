package main

import (
	"fmt"
	"net"
	"os"
	"sync"
)

var bufferPool = sync.Pool{
	New: func() any {
		fmt.Println("Pool making new buffer on heap")

		b := make([]byte, 1024)
		return &b
	},
}

func response(conn *net.UDPConn, addr *net.UDPAddr, data []byte) {
	message := string(data)
	response := []byte("Server got message: " + message)

	fmt.Println("Goroutine executed")

	_, err := conn.WriteToUDP(response, addr)
	if err != nil {
		fmt.Println("Error sending back data:", err)
	}
}

func main() {
	addr, err := net.ResolveUDPAddr("udp", ":9999")
	if err != nil {
		fmt.Println("Error resolving address:", err)
		os.Exit(1)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error listening UDP:", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("UDP Server running on  port :9999...")

	for {
		bufPtr := bufferPool.Get().(*[]byte)
		buf := *bufPtr

		n, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error reading data:", err)
			continue
		}

		// go response(conn, remoteAddr, buf[:n])

		go func(addr *net.UDPAddr, payload []byte, originalBuf *[]byte) {
			defer bufferPool.Put(originalBuf)
			response(conn, addr, payload)

		}(remoteAddr, buf[:n], bufPtr)

	}

}
