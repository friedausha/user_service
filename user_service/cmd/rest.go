package cmd

import (
	"bufio"
	"context"
	"fmt"
	handler2 "git.garena.com/frieda.hasanah/user_service/internal/handler"
	"git.garena.com/frieda.hasanah/user_service/utils/log"
	"github.com/spf13/cobra"
	"io"
	"net"
	"net/http"
	_ "net/http"
	_ "net/http/pprof"
	"os"
	"strings"
	"time"
)

var restCommand = &cobra.Command{
	Use:   "rest",
	Short: "Start TCP server",
	Run:   tcpServer,
}

func init() {
	rootCmd.AddCommand(restCommand)
}

var workerPoolSize = 50

func tcpServer(cmd *cobra.Command, args []string) {
	// srvAddress := MustHaveEnv("SERVER_ADDRESS")
	srvAddress := "0.0.0.0:8080"
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	// Start listening for TCP connections
	ln, err := net.Listen("tcp", srvAddress)
	if err != nil {
		log.Error(context.Background(), err, "Error starting TCP server")
		os.Exit(1)
	}
	defer ln.Close()

	fmt.Printf("TCP server started at %s\n", srvAddress)

	// Initialize handler
	h := handler2.NewHandler(authService, userService)
	connChannel := make(chan net.Conn, workerPoolSize)

	// Start worker pool
	// Cek whether it's thread safe atau ngga
	for i := 0; i < workerPoolSize; i++ {
		go worker(connChannel, h)
	}
	for {
		conn, err := ln.Accept()
		//fmt.Println("Accepted connection")
		if err != nil {
			log.Error(context.Background(), err, "Error accepting connection")
			continue
		}
		//go handleConnection(conn, h)
		connChannel <- conn
	}
}

func worker(connChannel chan net.Conn, h *handler2.Handler) {
	for conn := range connChannel {
		handleConnection(conn, h)
	}
}

func handleConnection(conn net.Conn, h *handler2.Handler) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for {
		// Set a read deadline to avoid hanging connections
		conn.SetReadDeadline(time.Now().Add(10 * time.Second))

		// Read message from client
		message, err := reader.ReadString('\n')
		if err != nil {
			//fmt.Println(err)
			// Handle different error types
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				log.Error(context.Background(), err, "Read timeout")
				break
			} else if err == io.EOF {
				//log.Info(context.Background(), err, "Client closed connection (EOF)")
				break
			} else {
				log.Error(context.Background(), err, "Unexpected error reading from connection")
				break
			}
		}

		// Process and respond to client
		response := h.HandleRequest(strings.TrimSpace(message))
		_, err = writer.WriteString(response + "\n")
		if err != nil {
			log.Error(context.Background(), err, "Error writing to connection")
			break
		}

		// Ensure data is flushed to the client
		err = writer.Flush()
		if err != nil {
			log.Error(context.Background(), err, "Error flushing writer")
			break
		}
	}
}
