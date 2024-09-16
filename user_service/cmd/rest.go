package cmd

import (
	"bufio"
	"context"
	"fmt"
	handler2 "git.garena.com/frieda.hasanah/user_service/internal/handler"
	"git.garena.com/frieda.hasanah/user_service/utils/log"
	"github.com/spf13/cobra"
	"net"
	"os"
	"strings"
)

var restCommand = &cobra.Command{
	Use:   "rest",
	Short: "Start TCP server",
	Run:   tcpServer,
}

func init() {
	rootCmd.AddCommand(restCommand)
}

func tcpServer(cmd *cobra.Command, args []string) {
	// srvAddress := MustHaveEnv("SERVER_ADDRESS")
	srvAddress := "0.0.0.0:8080"

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

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Error(context.Background(), err, "Error accepting connection")
			continue
		}

		go handleConnection(conn, h)
	}
}

func handleConnection(conn net.Conn, h *handler2.Handler) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Error(context.Background(), err, "Error reading from connection")
			break
		}

		response := h.HandleRequest(strings.TrimSpace(message))
		writer.WriteString(response + "\n")
		writer.Flush()
	}
}
