package cmd

import (
	"context"
	"fmt"
	handler2 "git.garena.com/frieda.hasanah/user_service/internal/handler"
	"git.garena.com/frieda.hasanah/user_service/utils/log"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
)

var restCommand = &cobra.Command{
	Use:   "rest",
	Short: "Start REST server",
	Run:   restServer,
}

func init() {
	rootCmd.AddCommand(restCommand)
}

// NewServer returns a new Echo server instance
func NewServer() (e *echo.Echo, g *echo.Group) {
	e = echo.New()
	g = e.Group("")
	return
}

func restServer(cmd *cobra.Command, args []string) {
	e, g := NewServer()

	g.GET("/healthcheck/liveness", func(c echo.Context) error {
		return c.String(200, "Calm down bro, I'm really-really healthy Bro!!!")
	})

	// Initialize handler
	handler2.Init(g, authService, userService)
	//srvAddress := MustHaveEnv("SERVER_ADDRESS")
	srvAddress := "0.0.0.0:8080"

	log.Error(context.Background(), nil, fmt.Sprintf("starting server at %s", srvAddress), e.Start(srvAddress))
}
