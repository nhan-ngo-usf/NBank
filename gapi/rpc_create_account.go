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

func (server *Server) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	// authenticate user
	payload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	// validate user input
	violations := validateCreateAccountRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	if payload.Username != req.Username {
		return nil, status.Errorf(codes.PermissionDenied, "unauthorized request")
	}

	arg := db.CreateAccountParams{
		UserName: 	payload.Username,
		Currency:	req.GetCurrency(),
	}
	
	if req.Balance != nil {
		arg.Balance = req.GetBalance()
	} else {
		arg.Balance = 0
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create account")
	}

	response := &pb.CreateAccountResponse{
		Account: convertAccount(account),
	}
	return response, nil
}

func validateCreateAccountRequest(req *pb.CreateAccountRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validate.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	if err := validate.ValidateCurrency(req.GetCurrency()); err != nil {
		violations = append(violations, fieldViolation("currency", err))
	}

	if err := validate.ValidateBalance(req.GetBalance()); err != nil {
		violations = append(violations, fieldViolation("balance", err))
	}

	return violations
}