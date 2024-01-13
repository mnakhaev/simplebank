package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/mnakhaev/simplebank/db/sqlc"
)

// Server serves HTTP requests for banking service.
type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer creates new HTTP server and setup routing.
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts/", server.listAccount)

	server.router = router
	return server
}

// Start runs the HTTP server on a specific address.
func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
