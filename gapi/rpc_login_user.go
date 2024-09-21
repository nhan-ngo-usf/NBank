package gapi

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	db "github.com/nhan-ngo-usf/NBank/db/sqlc"
	"github.com/nhan-ngo-usf/NBank/pb"
	"github.com/nhan-ngo-usf/NBank/util"
	"github.com/nhan-ngo-usf/NBank/validate"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	violations := validateLoginUserRequest(req)

	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	password := req.GetPassword()
	
	user, err := server.store.GetUser(ctx, req.GetUsername())

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to find user %s", err)
	}

	err = util.CheckPassword(password, user.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid password: %s", err)	
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.AccessTokenDuration,
	)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create access token")
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.RefreshTokenDuration,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create refresh token")
	}

	sessionID, err := uuid.NewRandom()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate session ID")
	}

	metadata := server.extractMetadata(ctx)
	session, err := server.store.CreateSession(ctx, db.CreateSessionParams {
		ID:		sessionID,
		Username:	user.Username,
		RefreshToken: 	refreshToken,
		UserAgent:	metadata.UserAgent,
		ClientIp:	metadata.ClientIP,
		IsBlocked:	false,
		ExpiresAt:	refreshPayload.ExpiredAt,
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session")
	}
	rsp := &pb.LoginUserResponse{
		User:		convertUser(user),
		SessionId:	session.ID.String(),
		AccessToken: 	accessToken,
		RefreshToken:   refreshToken,
		AccessTokenExpiresAt: timestamppb.New(accessPayload.ExpiredAt),
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
	}
	return rsp, nil
}

func validateLoginUserRequest(req *pb.LoginUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validate.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}	
	if err := validate.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	return violations
}