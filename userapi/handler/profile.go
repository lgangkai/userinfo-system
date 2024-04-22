package handler

import (
	"encoding/json"
	errs "errs"
	"github.com/asim/go-micro/v3/errors"
	"github.com/asim/go-micro/v3/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"protos/userinfo"
)

func (c *Client) GetProfile(context *gin.Context) {
	userId := c.getCookie(context, KEY_USER_ID)
	if userId == nil {
		return
	}

	r := userinfo.GetProfileRequest{
		UserId: userId.(uint64),
	}
	resp, err := c.UserinfoClient.GetProfile(context, &r)
	if err != nil {
		logger.Error("Call rpc server failed, error: ", err)
		code := errors.Parse(err.Error()).Code
		msg := errors.Parse(err.Error()).Detail
		context.JSON(http.StatusInternalServerError, gin.H{
			"code": code,
			"msg":  msg,
			"data": nil,
		})
		context.Abort()
		return
	}

	p, err := json.Marshal(resp.GetProfile())
	if err != nil {
		logger.Error("Marshal profile tp json failed, err: ", err.Error())
		context.JSON(http.StatusInternalServerError, gin.H{
			"code": errs.ERR_GET_PROFILE_FAILED,
			"msg":  errs.GetMsg(errs.ERR_GET_PROFILE_FAILED),
			"data": nil,
		})
		context.Abort()
		return
	}

	logger.Info("Handle get profile success, profile: ", string(p))
	context.JSON(http.StatusOK, gin.H{
		"code": errs.SUCCESS,
		"msg":  errs.GetMsg(errs.SUCCESS),
		"data": string(p),
	})
}

func (c *Client) UpdateProfile(context *gin.Context) {
	profile := &Profile{}
	if err := context.ShouldBind(profile); err != nil {
		logger.Error("Bind request data error, err: ", err.Error())
		context.JSON(http.StatusBadRequest, gin.H{
			"code": errs.ERR_UPDATE_PROFILE_FAILED,
			"msg":  errs.GetMsg(errs.ERR_UPDATE_PROFILE_FAILED),
			"data": nil,
		})
		context.Abort()
		return
	}

	userId := c.getCookie(context, KEY_USER_ID)
	if userId == nil {
		return
	}

	r := &userinfo.UpdateProfileRequest{
		Profile: &userinfo.Profile{
			UserId:   userId.(uint64),
			Username: profile.Username,
			Birthday: profile.Birthday,
		},
	}
	_, err := c.UserinfoClient.UpdateProfile(context, r)
	if err != nil {
		logger.Error("Call rpc server failed, error: ", err)
		code := errors.Parse(err.Error()).Code
		msg := errors.Parse(err.Error()).Detail
		context.JSON(http.StatusInternalServerError, gin.H{
			"code": code,
			"msg":  msg,
			"data": nil,
		})
		context.Abort()
		return
	}

	logger.Info("Handle update profile success.")
	context.JSON(http.StatusOK, gin.H{
		"code": errs.SUCCESS,
		"msg":  errs.GetMsg(errs.SUCCESS),
		"data": nil,
	})
}

func (c *Client) CreateProfile(context *gin.Context) {
	profile := &Profile{}
	if err := context.ShouldBind(profile); err != nil {
		logger.Error("Bind request data error, err: ", err.Error())
		context.JSON(http.StatusBadRequest, gin.H{
			"code": errs.ERR_CREATE_PROFILE_FAILED,
			"msg":  errs.GetMsg(errs.ERR_CREATE_PROFILE_FAILED),
			"data": nil,
		})
		context.Abort()
		return
	}

	userId := c.getCookie(context, KEY_USER_ID)
	email := c.getCookie(context, KEY_EMAIL)
	if userId == nil || email == nil {
		return
	}
	r := &userinfo.CreateProfileRequest{
		Profile: &userinfo.Profile{
			UserId:   userId.(uint64),
			Email:    email.(string),
			Username: profile.Username,
			Birthday: profile.Birthday,
		},
	}
	_, err := c.UserinfoClient.CreateProfile(context, r)
	if err != nil {
		logger.Error("Call rpc server failed, error: ", err)
		code := errors.Parse(err.Error()).Code
		msg := errors.Parse(err.Error()).Detail
		context.JSON(http.StatusInternalServerError, gin.H{
			"code": code,
			"msg":  msg,
			"data": nil,
		})
		context.Abort()
		return
	}

	logger.Info("Handle create profile success.")
	context.JSON(http.StatusOK, gin.H{
		"code": errs.SUCCESS,
		"msg":  errs.GetMsg(errs.SUCCESS),
		"data": nil,
	})
}

func (c *Client) getCookie(context *gin.Context, key string) any {
	userId, ok := context.Get(key)
	if !ok {
		logger.Error("Get cookie failed, key: ", key)
		context.JSON(http.StatusUnauthorized, gin.H{
			"code": errs.ERR_AUTH_FAILED,
			"msg":  errs.GetMsg(errs.ERR_AUTH_FAILED),
			"data": nil,
		})
		context.Abort()
		return nil
	}
	return userId
}
