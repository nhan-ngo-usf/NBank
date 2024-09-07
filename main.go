package main

import (
	"database/sql"
	"log"
	"net"

	_ "github.com/lib/pq"
	db "github.com/nhan-ngo-usf/NBank/db/sqlc"
	"github.com/nhan-ngo-usf/NBank/gapi"
	"github.com/nhan-ngo-usf/NBank/pb"
	"github.com/nhan-ngo-usf/NBank/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)


func main() {
    config, err := util.LoadConfig(".")
    if err != nil {
        log.Fatal("cannot load config:", err)
    }

    conn, err := sql.Open(config.DBDriver, config.DBSource)

    if err != nil {
        log.Fatal("cannot connect to db:", err)
    }

    store := db.NewStore(conn)
    server, err := gapi.NewServer(config, store)
    
    if err != nil {
        log.Fatal("cannot connect to server", err)
    }
    
    gRPCserver := grpc.NewServer()
    pb.RegisterBankServer(gRPCserver, server)
    reflection.Register(gRPCserver)

    listener, err := net.Listen("tcp", config.GRPCServerAddress)
    if err != nil {
        log.Fatal("Cannot create listener")
    }

    log.Printf("start gRPC server at %s", listener.Addr().String())
    err = gRPCserver.Serve(listener)
    if err != nil {
        log.Fatal("Cannot start gRPC server")
    }
}
