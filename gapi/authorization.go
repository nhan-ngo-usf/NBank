package gapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/nhan-ngo-usf/NBank/token"
	"google.golang.org/grpc/metadata"
)

const (
	authorizationHeader = "authorization"
	authorizationHeaderBearer = "bearer"
)
func (server *Server) authorizeUser(ctx context.Context) (*token.Payload, error){
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}

	values := md.Get(authorizationHeader)
	if len(values) == 0 {
		return nil, fmt.Errorf("missing authorization header")
	}
	
	authHeader := values[0]
	fields := strings.Fields(authHeader)
	if len(fields) < 2 {
		return nil, fmt.Errorf("invalid authorization header format")
	}

	authType := strings.ToLower(fields[0])
	if authType != authorizationHeaderBearer {
		return nil, fmt.Errorf("unsupported authorization method %s", authType)
	}

	authToken := fields[1]
	payload, err := server.tokenMaker.VerifyToken(authToken)
	if err != nil {
		return nil, fmt.Errorf("invalid token %s", err)
	}
	return payload, nil
}