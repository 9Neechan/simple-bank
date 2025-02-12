package api

import (
	"fmt"

	db "github.com/9Neechan/simple-bank/db/sqlc"
	"github.com/9Neechan/simple-bank/token"
	"github.com/9Neechan/simple-bank/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config     util.Config
	store      db.Store // *
	tokenMaker token.Maker
	router     *gin.Engine 
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	// tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users", server.createUser)      // http://localhost:8080/users
	router.POST("/users/login", server.loginUser) // http://localhost:8080/users/login

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRoutes.POST("/accounts", server.createAccount) // http://localhost:8080/accounts
	authRoutes.GET("/accounts/:id", server.getAccount) // http://localhost:8080/accounts/214
	authRoutes.GET("/accounts", server.listAccount)   // http://localhost:8080/accounts?page_id=1&page_size=5

	authRoutes.POST("/transfers", server.createTransfer) // http://localhost:8080/transfers

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
