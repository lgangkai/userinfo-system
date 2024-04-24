package profile

import (
	"context"
	"loggers"
	"protos/userinfo"
	"user-server/model"
	"user-server/service/profile"
)

type ProfileBiz struct {
	profileService *profile.ProfileService
	logger         *logger.Logger
}

func NewProfileBiz(profileService *profile.ProfileService, logger *logger.Logger) *ProfileBiz {
	return &ProfileBiz{
		profileService: profileService,
		logger:         logger,
	}
}

func (b *ProfileBiz) GetProfile(ctx context.Context, in *userinfo.GetProfileRequest, out *userinfo.GetProfileResponse) error {
	b.logger.Info(ctx, "Call ProfileBiz.GetProfile, request: ", in)
	id := in.GetUserId()
	p, err := b.profileService.GetProfile(ctx, id)
	if err != nil {
		b.logger.Error(ctx, "Get profile failed, err: ", err.Error())
		return err
	}
	out.Profile = &userinfo.Profile{
		Id:       p.Id,
		UserId:   p.UserId,
		Username: p.Username,
		Birthday: p.Birthday,
		Email:    p.Email,
		Avatar:   p.AvatarUrl,
	}
	b.logger.Info(ctx, "Call ProfileBiz.GetProfile successfully.")
	return nil
}

func (b *ProfileBiz) DeleteProfile(ctx context.Context, in *userinfo.DeleteProfileRequest, out *userinfo.DeleteProfileResponse) error {
	b.logger.Info(ctx, "Call ProfileBiz.DeleteProfile, request: ", in)
	id := in.GetUserId()
	err := b.profileService.DeleteProfile(ctx, id)
	if err != nil {
		b.logger.Error(ctx, "Delete profile failed, err: ", err.Error())
		return err
	}
	b.logger.Info(ctx, "Call ProfileBiz.DeleteProfile successfully.")
	return nil
}

func (b *ProfileBiz) CreateProfile(ctx context.Context, in *userinfo.CreateProfileRequest, out *userinfo.CreateProfileResponse) error {
	b.logger.Info(ctx, "Call ProfileBiz.CreateProfile, request: ", in)
	p := in.GetProfile()
	mp := &model.Profile{
		Id:        p.Id,
		UserId:    p.UserId,
		Username:  p.Username,
		Birthday:  p.Birthday,
		Email:     p.Email,
		AvatarUrl: p.Avatar,
	}
	err := b.profileService.CreateProfile(ctx, mp)
	if err != nil {
		b.logger.Error(ctx, "Create profile failed, err: ", err.Error())
		return err
	}
	b.logger.Info(ctx, "Call ProfileBiz.CreateProfile successfully.")
	return nil
}

func (b *ProfileBiz) UpdateProfile(ctx context.Context, in *userinfo.UpdateProfileRequest, out *userinfo.UpdateProfileResponse) error {
	b.logger.Info(ctx, "Call ProfileBiz.UpdateProfile, request: ", in)
	p := in.GetProfile()
	mp := &model.Profile{
		Id:        p.Id,
		UserId:    p.UserId,
		Username:  p.Username,
		Birthday:  p.Birthday,
		Email:     p.Email,
		AvatarUrl: p.Avatar,
	}
	err := b.profileService.UpdateProfile(ctx, p.UserId, mp)
	if err != nil {
		b.logger.Error(ctx, "Update profile failed, err: ", err.Error())
		return err
	}
	b.logger.Info(ctx, "Call ProfileBiz.UpdateProfile successfully.")
	return nil
}
