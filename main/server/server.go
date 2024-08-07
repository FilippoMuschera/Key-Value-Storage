package main

import (
	"SDCC/main/utils"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"strconv"
)

func main() {
	// Check command line arguments
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run server.go <server_index>")
		os.Exit(1)
	}
	index, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Invalid index")
		os.Exit(1)
	}

	// Check environment variable
	consistType := os.Getenv("CONSIST_TYPE")
	if consistType == "" {
		fmt.Println("Environment variable CONSIST_TYPE is not set")
		os.Exit(1)
	}
	fmt.Printf("CONSIST_TYPE: %s\n", consistType)

	// Set up RPC server
	sequential := NewKVSSequential()
	err = rpc.RegisterName("sequential", sequential)
	if err != nil {
		fmt.Println("Error registering RPC:", err)
		return
	}

	//local
	addr := "localhost:" + utils.GetServerPort(index)
	fmt.Println("Registering server ", index, " at address: ", addr)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}

	for {
		fmt.Println("Waiting for connection...")
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("SERVER: Errore nell'accettare la connessione dal client:", err)
			continue
		}

		// Avvia la gestione della connessione in un goroutine
		go func(conn net.Conn) {
			// Servi la connessione RPC
			fmt.Printf("Serving connection from %s\n", conn.RemoteAddr().String())
			rpc.ServeConn(conn)

			defer func() {
				err := conn.Close()
				if err != nil {
				}
			}()
		}(conn)

	}
}