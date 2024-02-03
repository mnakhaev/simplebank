package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/mnakhaev/simplebank/config"
	db "github.com/mnakhaev/simplebank/db/sqlc"
	"github.com/mnakhaev/simplebank/token"
)

// Server serves HTTP requests for banking service.
type Server struct {
	config     config.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

// NewServer creates new HTTP server and setup routing.
func NewServer(config config.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	// RegisterValidation is struct method of `*validator.Validate`, thus need to cast.
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// all currencies can be replaced with `currency` in validator.
		_ = v.RegisterValidation("currency", validCurrency)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}
	server.setupRouter()

	return server, nil
}

func (s *Server) setupRouter() {
	router := gin.Default()
	router.POST("/users", s.createUser)
	router.POST("/users/login", s.loginUser)

	authRoutes := router.Group("/").Use(authMiddleware(s.tokenMaker))
	// all other routes (which need to be protected) should be added to authRoutes group
	authRoutes.POST("/accounts", s.createAccount)
	authRoutes.GET("/accounts/:id", s.getAccount)
	authRoutes.GET("/accounts/", s.listAccount)

	authRoutes.POST("/transfers", s.createTransfer)

	s.router = router
}

// Start runs the HTTP server on a specific address.
func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
