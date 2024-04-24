package account

import (
	"context"
	"loggers"
	"protos/userinfo"
	"user-server/service/account"
)

type AccountBiz struct {
	accountService *account.AccountService
	logger         *logger.Logger
}

func NewAccountBiz(accountService *account.AccountService, logger *logger.Logger) *AccountBiz {
	return &AccountBiz{
		accountService: accountService,
		logger:         logger,
	}
}

func (b *AccountBiz) Register(ctx context.Context, in *userinfo.RegisterRequest, out *userinfo.RegisterResponse) error {
	b.logger.Info(ctx, "Call AccountBiz.Register, request: ", in)
	err := b.accountService.Register(ctx, in.GetEmail(), in.GetPassword())
	if err != nil {
		b.logger.Error(ctx, "Register failed, err: ", err.Error())
		return err
	}
	b.logger.Info(ctx, "Call AccountBiz.Register successfully.")
	return nil
}

func (b *AccountBiz) Login(ctx context.Context, in *userinfo.LoginRequest, out *userinfo.LoginResponse) error {
	b.logger.Info(ctx, "Call AccountBiz.Login, request: ", in)
	token, err := b.accountService.Login(ctx, in.GetEmail(), in.GetPassword())
	if err != nil {
		b.logger.Error(ctx, "Login failed, err: ", err.Error())
		return err
	}
	b.logger.Info(ctx, "Call AccountBiz.Login successfully.")
	out.Token = *token
	return nil
}

// Logout Currently no implement for logout in service.
// Just remove token in api gateway.
func (b *AccountBiz) Logout(ctx context.Context, in *userinfo.LogoutRequest, out *userinfo.LogoutResponse) error {
	return nil
}

func (b *AccountBiz) Authenticate(ctx context.Context, in *userinfo.AuthRequest, out *userinfo.AuthResponse) error {
	b.logger.Info(ctx, "Call AccountBiz.Authenticate, request: ", in)
	userId, email, err := b.accountService.Authenticate(ctx, in.GetToken())
	if err != nil {
		b.logger.Error(ctx, "Authenticate failed, err: ", err.Error())
		return err
	}
	b.logger.Info(ctx, "Call AccountBiz.Authenticate successfully.")
	out.UserId = userId
	out.Email = email
	return nil
}
