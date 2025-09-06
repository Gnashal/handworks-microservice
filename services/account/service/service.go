package service

import (
	"context"
	"handworks/common/grpc/account"
	"handworks/common/utils"
)

type AccountService struct {
	L    *utils.Logger
	Grpc account.AccountServiceServer
	account.UnimplementedAccountServiceServer
}

func (a *AccountService) SignUp(c context.Context, in *account.SignUpRequest) (*account.SignUpResponse, error) {

	return nil, nil
}
