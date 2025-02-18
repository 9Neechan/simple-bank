package grpcapi

import (
	"fmt"

	db "github.com/9Neechan/simple-bank/db/sqlc"
	"github.com/9Neechan/simple-bank/pb"
	"github.com/9Neechan/simple-bank/token"
	"github.com/9Neechan/simple-bank/util"
)

// Server обслуживает gRPC запросы нашего банковского сервиса.
type Server struct {
	pb.UnimplementedSimpleBankServer // сервер уже сможет принимать RPC вызовы до того, как они будут фактически реализованы
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}

// NewServer создаёт новый GRPc сервер и настраивает маршрутизацию.
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	return server, nil
}
