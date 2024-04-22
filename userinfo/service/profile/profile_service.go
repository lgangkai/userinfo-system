package profile

import (
	"context"
	errs "errs"
	"github.com/asim/go-micro/v3/logger"
	"user-server/dao"
	"user-server/model"
)

type ProfileService struct {
	profileDao *dao.ProfileDao
}

func NewProfileService(profileDao *dao.ProfileDao) *ProfileService {
	return &ProfileService{
		profileDao: profileDao,
	}
}

func (s *ProfileService) GetProfile(ctx context.Context, userId uint64) (*model.Profile, error) {
	logger.Info("Call ProfileService.GetProfile, user_id: ", userId)
	profile, err := s.profileDao.GetProfileById(ctx, userId)
	if err != nil {
		logger.Error("Fail to get profile, err:", err.Error())
		return nil, errs.New(errs.ERR_GET_PROFILE_FAILED)
	}
	return profile, nil
}

func (s *ProfileService) UpdateProfile(ctx context.Context, userId uint64, profile *model.Profile) error {
	logger.Info("Call ProfileService.UpdateProfile, user_id: ", userId)
	err := s.profileDao.Update(ctx, userId, profile)
	if err != nil {
		logger.Error("Fail to update profile, err:", err.Error())
		return errs.New(errs.ERR_UPDATE_PROFILE_FAILED)
	}
	return nil
}

func (s *ProfileService) DeleteProfile(ctx context.Context, userId uint64) error {
	logger.Info("Call ProfileService.DeleteProfile, user_id: ", userId)
	err := s.profileDao.Delete(ctx, userId)
	if err != nil {
		logger.Error("Fail to delete profile, err:", err.Error())
		return errs.New(errs.ERR_DELETE_PROFILE_FAILED)
	}
	return nil
}

func (s *ProfileService) CreateProfile(profile *model.Profile) error {
	logger.Info("Call ProfileService.CreateProfile, profile: ", profile)
	err := s.profileDao.Insert(profile)
	if err != nil {
		logger.Error("Fail to delete profile, err:", err.Error())
		return errs.New(errs.ERR_CREATE_PROFILE_FAILED)
	}
	return nil
}
