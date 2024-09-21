package gapi

import (
	db "github.com/nhan-ngo-usf/NBank/db/sqlc"
	"github.com/nhan-ngo-usf/NBank/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertUser(user db.User) *pb.User {
	return &pb.User{
		Username:	user.Username,
		FullName:	user.FullName,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt:	timestamppb.New(user.CreatedAt),
	}
}

func convertAccount(account db.Account) *pb.Account {
	return &pb.Account {
		Username: account.UserName,
		Balance: account.Balance,
		Currency: account.Currency,
		CreatedAt: timestamppb.New(account.CreatedAt),
	}
}

func convertListAccount(accounts []db.Account) ([]*pb.Account){
	var pbAccounts []*pb.Account
	for _, account := range accounts {
		pbAccounts = append(pbAccounts, &pb.Account{
			Username: account.UserName,
			Balance: account.Balance,
			Currency: account.Currency,
			CreatedAt: timestamppb.New(account.CreatedAt),
		})
	}
	return pbAccounts
}