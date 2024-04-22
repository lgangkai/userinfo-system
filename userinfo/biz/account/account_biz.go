package account

import (
	"context"
	"github.com/asim/go-micro/v3/logger"
	"protos/userinfo"
	"user-server/service/account"
)

type AccountBiz struct {
	accountService *account.AccountService
}

func NewAccountBiz(accountService *account.AccountService) *AccountBiz {
	return &AccountBiz{
		accountService: accountService,
	}
}

func (b *AccountBiz) Register(ctx context.Context, in *userinfo.RegisterRequest, out *userinfo.RegisterResponse) error {
	logger.Info("Call AccountBiz.Register, request: ", in)
	err := b.accountService.Register(in.GetEmail(), in.GetPassword())
	if err != nil {
		logger.Error("Register failed, err: ", err.Error())
		return err
	}
	logger.Info("Call AccountBiz.Register successfully.")
	return nil
}

func (b *AccountBiz) Login(ctx context.Context, in *userinfo.LoginRequest, out *userinfo.LoginResponse) error {
	logger.Info("Call AccountBiz.Login, request: ", in)
	token, err := b.accountService.Login(in.GetEmail(), in.GetPassword())
	if err != nil {
		logger.Error("Login failed, err: ", err.Error())
		return err
	}
	logger.Info("Call AccountBiz.Login successfully.")
	out.Token = *token
	return nil
}

// Logout Currently no implement for logout in service.
// Just remove token in api gateway.
func (b *AccountBiz) Logout(ctx context.Context, in *userinfo.LogoutRequest, out *userinfo.LogoutResponse) error {
	return nil
}

func (b *AccountBiz) Authenticate(ctx context.Context, in *userinfo.AuthRequest, out *userinfo.AuthResponse) error {
	logger.Info("Call AccountBiz.Authenticate, request: ", in)
	userId, email, err := b.accountService.Authenticate(in.GetToken())
	if err != nil {
		logger.Error("Authenticate failed, err: ", err.Error())
		return err
	}
	logger.Info("Call AccountBiz.Authenticate successfully.")
	out.UserId = userId
	out.Email = email
	return nil
}
