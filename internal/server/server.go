package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/leonardoberlatto/go-url-shortener/internal/routes"
)

type Server struct {
	router *gin.Engine
}

func Initialize() *Server {
	router := gin.Default()
	server := &Server{router: router}
	routes.SetupRoutes(server.router)

	return server
}

func (server *Server) Start(port string) error {
	return server.router.Run(fmt.Sprintf(":%s", port))
}
