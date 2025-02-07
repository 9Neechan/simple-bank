package api

import (
	db "github.com/9Neechan/simple-bank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  db.Store // *
	router *gin.Engine
}

func NewServer(store db.Store) *Server {  // *
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.createAccount) // http://localhost:8080/accounts
	router.GET("/accounts/:id", server.getAccount) // http://localhost:8080/accounts/214
	router.GET("/accounts", server.listAccounts) // http://localhost:8080/accounts?page_id=1&page_size=5 

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
