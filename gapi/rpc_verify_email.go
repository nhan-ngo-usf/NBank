package gapi

import (
	"context"
	"fmt"

	db "github.com/nhan-ngo-usf/NBank/db/sqlc"
	"github.com/nhan-ngo-usf/NBank/pb"
	"github.com/nhan-ngo-usf/NBank/validate"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


func (server *Server)VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error){
	violations := validateVerifyEmailRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	txResult, err := server.store.VerifyEmailTx(ctx, db.VerifyEmailTxParams{
		EmailID:    req.GetEmailId(),
		VerifyCode: req.GetVerifyCode(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to verify email")
	}

	rsp := &pb.VerifyEmailResponse{
		IsVerified: txResult.User.IsEmailVerified,
	}
	return rsp, nil
}

func validateVerifyEmailRequest(req *pb.VerifyEmailRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validate.ValidateEmailID(req.GetEmailId()); err != nil {
		violations = append(violations, fieldViolation("email_id", err))
	}
	fmt.Println(req.GetVerifyCode())
	if err := validate.ValidateVerifyCode(req.GetVerifyCode()); err != nil {
		violations = append(violations, fieldViolation("verify_code", err))
	}

	return violations
}