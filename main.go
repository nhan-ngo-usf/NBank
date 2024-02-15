package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	db "github.com/nhan-ngo-usf/NBank/db/sqlc"
	"github.com/nhan-ngo-usf/NBank/gapi"
	"github.com/nhan-ngo-usf/NBank/pb"
	"github.com/nhan-ngo-usf/NBank/util"
	"google.golang.org/grpc"
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
        log.Fatal("cannot connect to server")
    }
    
    gRPCserver := grpc.NewServer()
    pb.RegisterBankServer(gRPCserver, server)
    
}
