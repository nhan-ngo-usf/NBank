package main

import (
	"context"
	"database/sql"
	"net"
	"net/http"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hibiken/asynq"
	_ "github.com/lib/pq"
	db "github.com/nhan-ngo-usf/NBank/db/sqlc"
	_ "github.com/nhan-ngo-usf/NBank/doc/statik"
	"github.com/nhan-ngo-usf/NBank/gapi"
	"github.com/nhan-ngo-usf/NBank/mail"
	"github.com/nhan-ngo-usf/NBank/pb"
	"github.com/nhan-ngo-usf/NBank/util"
	"github.com/nhan-ngo-usf/NBank/worker"
	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)


func main() {
    config, err := util.LoadConfig(".")
    if err != nil {
        log.Fatal().Msg("cannot load config:")
    }
    if config.Environment == "development"{
        log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
    }
    conn, err := sql.Open(config.DBDriver, config.DBSource)

    if err != nil {
        log.Fatal().Msg("cannot connect to db:")
    }

    store := db.NewStore(conn)
    
    redisOpt := asynq.RedisClientOpt{
        Addr: config.RedisAddress,
    }

    taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

    go RunTaskProcessor(config, redisOpt, store)
    go RunGatewayServer(config, store, taskDistributor)
    RunGRPCServer(config, store, taskDistributor)
}

func RunGRPCServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
    server, err := gapi.NewServer(config, store, taskDistributor)  
    if err != nil {
        log.Fatal().Msg("cannot connect to server")
    }
    
    grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
    gRPCserver := grpc.NewServer(grpcLogger)
    pb.RegisterBankServer(gRPCserver, server)
    reflection.Register(gRPCserver)

    listener, err := net.Listen("tcp", config.GRPCServerAddress)
    if err != nil {
        log.Fatal().Msg("Cannot create listener")
    }

    log.Info().Msgf("start gRPC server at %s", listener.Addr().String())
    err = gRPCserver.Serve(listener)
    if err != nil {
        log.Fatal().Msg("Cannot start gRPC server")
    }
}

func RunGatewayServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
    server, err := gapi.NewServer(config, store, taskDistributor)  
    if err != nil {
        log.Fatal().Msg("cannot connect to server")
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
        log.Fatal().Msg("cannot register handler")
    }
    mux := http.NewServeMux()
    mux.Handle("/", grpcMux)

    statikFS, err := fs.New()
    if err != nil {
        log.Fatal().Msg("cannot create statik fs")
    }

    swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
    mux.Handle("/swagger/", swaggerHandler)

    listener, err := net.Listen("tcp", config.HTTPServerAddress)
    if err != nil {
        log.Fatal().Msg("Cannot create listener")
    }

    log.Info().Msgf("start HTTP gateway at %s", listener.Addr().String())
    handler := gapi.HttpLogger(mux)
    err = http.Serve(listener, handler)
    if err != nil {
        log.Fatal().Msg("Cannot start HTTP gateway server")
    }
}

func RunTaskProcessor(// ctx context.Context,
    config util.Config,
	redisOpt asynq.RedisClientOpt,
	store db.Store,
) {
	mailer := mail.NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store, mailer)

	log.Info().Msg("start task processor")
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start task processor")
	}

	
}
