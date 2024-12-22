package gapi

import (
	"fmt"

	db "github.com/nhan-ngo-usf/NBank/db/sqlc"
	"github.com/nhan-ngo-usf/NBank/pb"
	"github.com/nhan-ngo-usf/NBank/token"
	"github.com/nhan-ngo-usf/NBank/util"
)



type Server struct {
	pb.UnimplementedBankServer
	store	 	db.Store
	config	 	util.Config
	tokenMaker 	token.Maker
}
func NewServer(config util.Config, store db.Store) (*Server, error){
	
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		store: store,
		config: config,
		tokenMaker: tokenMaker,
	}

	return server, nil
}

