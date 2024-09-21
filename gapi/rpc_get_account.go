package gapi

import (
	"context"
	"strconv"

	"github.com/nhan-ngo-usf/NBank/pb"
	"github.com/nhan-ngo-usf/NBank/validate"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server)GetAccount(ctx context.Context, req *pb.GetAccountRequest) (*pb.GetAccountResponse, error) {
	// authenticate user
	payload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}
	
	// validate user input
	violations := validateGetAccountRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	id, err := strconv.Atoi(req.AccountId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot convert id to int %s", err)
	}

	account, err := server.store.GetAccount(ctx, int64(id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create account %s", err)
	}
	
	if account.UserName != payload.Username {
		return nil, status.Errorf(codes.PermissionDenied, "account belongs to different user")
	}

	response := &pb.GetAccountResponse{
		Account: convertAccount(account),
	}
	return response, nil
}

func validateGetAccountRequest(req *pb.GetAccountRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	
	if err := validate.ValidateAccountID(req.GetAccountId()); err != nil {
		violations = append(violations, fieldViolation("account_id", err))
	}
	return violations
}