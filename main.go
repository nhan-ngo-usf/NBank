package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	db "github.com/nhan-ngo-usf/NBank/db/sqlc"
	"github.com/nhan-ngo-usf/NBank/gapi"
	"github.com/nhan-ngo-usf/NBank/pb"
	"github.com/nhan-ngo-usf/NBank/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
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
    go RunGatewayServer(config, store)
    RunGRPCServer(config, store)
}
func RunGRPCServer(config util.Config, store db.Store) {
    server, err := gapi.NewServer(config, store)  
    if err != nil {
        log.Fatal("cannot connect to server", err)
    }
    
    gRPCserver := grpc.NewServer()
    pb.RegisterBankServer(gRPCserver, server)
    reflection.Register(gRPCserver)

    listener, err := net.Listen("tcp", config.GRPCServerAddress)
    if err != nil {
        log.Fatal("Cannot create listener", err)
    }

    log.Printf("start gRPC server at %s", listener.Addr().String())
    err = gRPCserver.Serve(listener)
    if err != nil {
        log.Fatal("Cannot start gRPC server", err)
    }
}

func RunGatewayServer(config util.Config, store db.Store) {
    server, err := gapi.NewServer(config, store)  
    if err != nil {
        log.Fatal("cannot connect to server", err)
    }
    
    jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
            MarshalOptions: protojson.MarshalOptions{
                UseProtoNames: true,
            },
            UnmarshalOptions: protojson.UnmarshalOptions{
                DiscardUnknown: true,
            },
        })

    grpcMux := runtime.NewServeMux(jsonOption)
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    err = pb.RegisterBankHandlerServer(ctx, grpcMux, server)
    if err != nil {
        log.Fatal("cannot register handler", err)
    }
    mux := http.NewServeMux()
    mux.Handle("/", grpcMux)

    listener, err := net.Listen("tcp", config.HTTPServerAddress)
    if err != nil {
        log.Fatal("Cannot create listener", err)
    }

    log.Printf("start HTTP gateway at %s", listener.Addr().String())
    err = http.Serve(listener, mux)
    if err != nil {
        log.Fatal("Cannot start HTTP gateway server")
    }
}