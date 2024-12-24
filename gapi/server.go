package gapi

import (
	"fmt"

	db "github.com/nhan-ngo-usf/NBank/db/sqlc"
	"github.com/nhan-ngo-usf/NBank/pb"
	"github.com/nhan-ngo-usf/NBank/token"
	"github.com/nhan-ngo-usf/NBank/util"
	"github.com/nhan-ngo-usf/NBank/worker"
)



type Server struct {
	pb.UnimplementedBankServer
	store	 	db.Store
	config	 	util.Config
	tokenMaker 	token.Maker
	TaskDistributor worker.TaskDistributor
}
func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error){
	
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		store: store,
		config: config,
		tokenMaker: tokenMaker,
		TaskDistributor: taskDistributor,
	}
	return server, nil
}

