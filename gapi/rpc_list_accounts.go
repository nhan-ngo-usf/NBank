package gapi

import (
	"context"

	db "github.com/nhan-ngo-usf/NBank/db/sqlc"
	"github.com/nhan-ngo-usf/NBank/pb"
	"github.com/nhan-ngo-usf/NBank/validate"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) ListAccount(ctx context.Context, req *pb.ListAccountsRequest) (*pb.ListAccountsResponse, error){
	// validate user inputs
	violations := validateListAccountRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	// authenticate users
	payload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	if payload.Username != req.GetUsername(){
		return nil, status.Errorf(codes.PermissionDenied, "cannot access other users' accounts")
	}
	
	arg := db.ListAccountsParams{
		UserName: 	payload.Username,
		Limit: 		req.GetLimit(),
		Offset: 	req.GetOffset(),
	}

	accounts, err := server.store.ListAccounts(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot get list of accounts %s", err)
	}

	request := &pb.ListAccountsResponse{
		Accounts: convertListAccount(accounts),
	}
	return request, nil
}


func validateListAccountRequest(req *pb.ListAccountsRequest)(violations []*errdetails.BadRequest_FieldViolation) {
	if err := validate.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	return violations
}