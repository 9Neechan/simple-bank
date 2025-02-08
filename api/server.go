package api

import (
	db "github.com/9Neechan/simple-bank/db/sqlc"
	
	"github.com/go-playground/validator/v10"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type Server struct {
	store  db.Store // *
	router *gin.Engine
}

func NewServer(store db.Store) *Server { // *
	server := &Server{store: store}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	router.POST("/accounts", server.createAccount)   // http://localhost:8080/accounts
	router.GET("/accounts/:id", server.getAccount)   // http://localhost:8080/accounts/214
	router.GET("/accounts", server.listAccounts)     // http://localhost:8080/accounts?page_id=1&page_size=5
	router.POST("/transfers", server.createTransfer) // http://localhost:8080/transfers

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
