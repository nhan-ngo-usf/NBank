package gapi

import (
	"context"
	"database/sql"
	"time"

	db "github.com/nhan-ngo-usf/NBank/db/sqlc"
	"github.com/nhan-ngo-usf/NBank/pb"
	"github.com/nhan-ngo-usf/NBank/util"
	"github.com/nhan-ngo-usf/NBank/validate"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	payload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	violations := validateUpdateUserRequest(req)
	if violations != nil {

		return nil, invalidArgumentError(violations)
	}

	if payload.Username != req.GetUsername() {
		return nil, status.Errorf(codes.PermissionDenied, "cannot update other user's info")
	}
	arg := db.UpdateUserParams {
		Username:	req.GetUsername(),
		FullName:	sql.NullString{
			String: req.GetFullName(),
			Valid: req.FullName != nil && req.GetFullName() != "",
		},
		Email:	sql.NullString{
			String: req.GetEmail(),
			Valid: req.Email != nil && req.GetEmail() != "",
		},
	}

	if req.Password != nil && req.GetPassword() != ""{
		hashedPassword, err := util.HashPassword(req.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
			
		}

		arg.HashedPassword = sql.NullString{
			String: hashedPassword,
			Valid: true,
		}

		arg.PasswordChangedAt = sql.NullTime{
			Time: time.Now(),
			Valid: true,
		}
	}

	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows{
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "fail to update user: %s", err)
	}
	
	response := &pb.UpdateUserResponse{
		User: convertUser(user),
	}
	return response, nil
}

func validateUpdateUserRequest(req *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validate.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if req.FullName != nil && req.GetFullName() != "" {
		if err := validate.ValidateFullName(req.GetFullName()); err != nil {
			violations = append(violations, fieldViolation("full_name", err))
		}
	}

	if req.Email != nil && req.GetEmail() != ""{
		if err := validate.ValidateEmail(req.GetEmail()); err != nil {
			violations = append(violations, fieldViolation("email", err))
		}
	}

	if req.Password != nil && req.GetPassword() != ""{
		if err := validate.ValidatePassword(req.GetPassword()); err != nil {
			violations = append(violations, fieldViolation("password", err))
		}
	}
	return violations
}