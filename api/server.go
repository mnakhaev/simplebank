package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
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

	// RegisterValidation is struct method of `*validator.Validate`, thus need to cast.
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// all currencies can be replaced with `currency` in validator.
		_ = v.RegisterValidation("currency", validCurrency)
	}

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts/", server.listAccount)
	router.POST("/transfers", server.createTransfer)

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