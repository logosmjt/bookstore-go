package gapi

import (
	"fmt"

	db "github.com/logosmjt/bookstore-go/db/sqlc"
	"github.com/logosmjt/bookstore-go/pb"
	"github.com/logosmjt/bookstore-go/token"
	"github.com/logosmjt/bookstore-go/util"
)

type Server struct {
	pb.UnimplementedBookStoreServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}

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
