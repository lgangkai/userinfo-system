package profile

import (
	"context"
	"github.com/asim/go-micro/v3/logger"
	"protos/userinfo"
	"user-server/model"
	"user-server/service/profile"
)

type ProfileBiz struct {
	profileService *profile.ProfileService
}

func NewProfileBiz(profileService *profile.ProfileService) *ProfileBiz {
	return &ProfileBiz{profileService: profileService}
}

func (b *ProfileBiz) GetProfile(ctx context.Context, in *userinfo.GetProfileRequest, out *userinfo.GetProfileResponse) error {
	logger.Info("Call ProfileBiz.GetProfile, request: ", in)
	id := in.GetUserId()
	p, err := b.profileService.GetProfile(ctx, id)
	if err != nil {
		logger.Error("Get profile failed, err: ", err.Error())
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
	logger.Info("Call ProfileBiz.GetProfile successfully.")
	return nil
}

func (b *ProfileBiz) DeleteProfile(ctx context.Context, in *userinfo.DeleteProfileRequest, out *userinfo.DeleteProfileResponse) error {
	logger.Info("Call ProfileBiz.DeleteProfile, request: ", in)
	id := in.GetUserId()
	err := b.profileService.DeleteProfile(ctx, id)
	if err != nil {
		logger.Error("Delete profile failed, err: ", err.Error())
		return err
	}
	logger.Info("Call ProfileBiz.DeleteProfile successfully.")
	return nil
}

func (b *ProfileBiz) CreateProfile(ctx context.Context, in *userinfo.CreateProfileRequest, out *userinfo.CreateProfileResponse) error {
	logger.Info("Call ProfileBiz.CreateProfile, request: ", in)
	p := in.GetProfile()
	mp := &model.Profile{
		Id:        p.Id,
		UserId:    p.UserId,
		Username:  p.Username,
		Birthday:  p.Birthday,
		Email:     p.Email,
		AvatarUrl: p.Avatar,
	}
	err := b.profileService.CreateProfile(mp)
	if err != nil {
		logger.Error("Create profile failed, err: ", err.Error())
		return err
	}
	logger.Info("Call ProfileBiz.CreateProfile successfully.")
	return nil
}

func (b *ProfileBiz) UpdateProfile(ctx context.Context, in *userinfo.UpdateProfileRequest, out *userinfo.UpdateProfileResponse) error {
	logger.Info("Call ProfileBiz.UpdateProfile, request: ", in)
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
		logger.Error("Update profile failed, err: ", err.Error())
		return err
	}
	logger.Info("Call ProfileBiz.UpdateProfile successfully.")
	return nil
}
