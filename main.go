package main

import (
	"flag"
	"io"
	"log"
	"net"
)

var (
	listenAddr = flag.String("listen-addr", ":8080", "server listen address")

	// List of backend servers
	server = []string{
		"localhost:5001",
		"localhost:5002",
		"localhost:5003",
	}
)

func main() {
	// Parse command line flags
	flag.Parse()

	// Listen for incoming connections
	listener, err := net.Listen("tcp", *listenAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer listener.Close()

	// Accept incoming connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %v", err)
			continue
		}

		// Choose a backend server
		serverAddr := chooseServer()

		// Start proxying data between client and server
		go proxy(serverAddr, conn)
	}
}

// proxy handles proxying data between client and server
func proxy(serverAddr string, conn net.Conn) {
	// Connect to the chosen backend server
	sr, err := net.Dial("tcp", serverAddr)
	if err != nil {
		log.Printf("failed to dial server %s: %v", serverAddr, err)
		return
	}
	defer sr.Close()

	// Copy data from client to server
	go io.Copy(sr, conn)

	// Copy data from server to client
	go io.Copy(conn, sr)
}

// chooseServer selects a backend server to forward the request
func chooseServer() string {
	// For now, always choose the first server in the list
	return server[0]
}
