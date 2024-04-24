package handler

import (
	"context"
	"loggers"
	"protos/userinfo"
	"user-server/biz/account"
	"user-server/biz/profile"
)

type UserinfoHandlerImpl struct {
	accountBiz *account.AccountBiz
	profileBiz *profile.ProfileBiz
}

func NewUserinfoHandlerImpl(profileBiz *profile.ProfileBiz, accountBiz *account.AccountBiz) *UserinfoHandlerImpl {
	return &UserinfoHandlerImpl{
		accountBiz: accountBiz,
		profileBiz: profileBiz,
	}
}

func (h *UserinfoHandlerImpl) GetProfile(ctx context.Context, in *userinfo.GetProfileRequest, out *userinfo.GetProfileResponse) error {
	return h.profileBiz.GetProfile(getTraceContext(ctx, in.GetRequestId(), in.GetUserId()), in, out)
}

func (h *UserinfoHandlerImpl) DeleteProfile(ctx context.Context, in *userinfo.DeleteProfileRequest, out *userinfo.DeleteProfileResponse) error {
	return h.profileBiz.DeleteProfile(getTraceContext(ctx, in.GetRequestId(), in.GetUserId()), in, out)
}

func (h *UserinfoHandlerImpl) CreateProfile(ctx context.Context, in *userinfo.CreateProfileRequest, out *userinfo.CreateProfileResponse) error {
	return h.profileBiz.CreateProfile(getTraceContext(ctx, in.GetRequestId(), in.GetProfile().GetUserId()), in, out)
}

func (h *UserinfoHandlerImpl) UpdateProfile(ctx context.Context, in *userinfo.UpdateProfileRequest, out *userinfo.UpdateProfileResponse) error {
	return h.profileBiz.UpdateProfile(getTraceContext(ctx, in.GetRequestId(), in.GetProfile().GetUserId()), in, out)
}

func (h *UserinfoHandlerImpl) Register(ctx context.Context, in *userinfo.RegisterRequest, out *userinfo.RegisterResponse) error {
	return h.accountBiz.Register(getTraceContext(ctx, in.GetRequestId(), 0), in, out)
}

func (h *UserinfoHandlerImpl) Login(ctx context.Context, in *userinfo.LoginRequest, out *userinfo.LoginResponse) error {
	return h.accountBiz.Login(getTraceContext(ctx, in.GetRequestId(), 0), in, out)
}

func (h *UserinfoHandlerImpl) Logout(ctx context.Context, in *userinfo.LogoutRequest, out *userinfo.LogoutResponse) error {
	return h.accountBiz.Logout(getTraceContext(ctx, in.GetRequestId(), 0), in, out)
}

func (h *UserinfoHandlerImpl) Authenticate(ctx context.Context, in *userinfo.AuthRequest, out *userinfo.AuthResponse) error {
	return h.accountBiz.Authenticate(getTraceContext(ctx, in.GetRequestId(), 0), in, out)
}

func getTraceContext(ctx context.Context, requestId string, userId uint64) context.Context {
	return context.WithValue(ctx, logger.TraceDataKey{}, logger.TraceData{
		RequestId: requestId,
		UserId:    userId,
	})
}
