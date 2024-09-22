package cmd

import (
	"bufio"
	"context"
	"fmt"
	"git.garena.com/frieda.hasanah/user_service/internal/handler"
	"git.garena.com/frieda.hasanah/user_service/utils/log"
	"github.com/spf13/cobra"
	"io"
	"net"
	"net/http"
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

var workerPoolSize = 500
var connTimeout = 10 * time.Second

func tcpServer(cmd *cobra.Command, args []string) {
	srvAddress := "0.0.0.0:8080"
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	ln, err := net.Listen("tcp", srvAddress)
	if err != nil {
		log.Error(context.Background(), err, "Error starting TCP server")
		os.Exit(1)
	}
	defer ln.Close()

	fmt.Printf("TCP server started at %s\n", srvAddress)

	h := handler.NewHandler(authService, userService)
	connChannel := make(chan net.Conn, workerPoolSize)

	for i := 0; i < workerPoolSize; i++ {
		go worker(connChannel, h)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Error(context.Background(), err, "Error accepting connection")
			continue
		}
		connChannel <- conn
		//go handleConnection(conn, h)
	}
}

func worker(connChannel chan net.Conn, h *handler.Handler) {
	for conn := range connChannel {
		handleConnection(conn, h)
	}
}

func handleConnection(conn net.Conn, h *handler.Handler) {
	defer conn.Close()

	// Set a read deadline to avoid hanging connections
	conn.SetDeadline(time.Now().Add(connTimeout))

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				log.Error(context.Background(), err, "Read timeout")
				break
			} else if err == io.EOF {
				break
			} else {
				log.Error(context.Background(), err, "Unexpected error reading from connection")
				break
			}
		}

		response := h.HandleRequest(strings.TrimSpace(message))
		_, err = writer.WriteString(response + "\n")
		if err != nil {
			log.Error(context.Background(), err, "Error writing to connection")
			break
		}

		err = writer.Flush()
		if err != nil {
			log.Error(context.Background(), err, "Error flushing writer")
			break
		}
	}
}
