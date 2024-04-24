package handler

import (
	"encoding/json"
	errs "errs"
	"github.com/asim/go-micro/v3/errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"protos/userinfo"
)

func (c *Client) GetProfile(context *gin.Context) {
	userId := c.getAuthedData(context, KEY_USER_ID)
	if userId == nil {
		return
	}

	r := userinfo.GetProfileRequest{
		UserId:    userId.(uint64),
		RequestId: GetRequestId(context),
	}
	resp, err := c.userinfoClient.GetProfile(context, &r)
	if err != nil {
		c.logger.Error(c.context, "Call rpc server failed, error: ", err)
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
		c.logger.Error(c.context, "Marshal profile tp json failed, err: ", err.Error())
		context.JSON(http.StatusInternalServerError, gin.H{
			"code": errs.ERR_GET_PROFILE_FAILED,
			"msg":  errs.GetMsg(errs.ERR_GET_PROFILE_FAILED),
			"data": nil,
		})
		context.Abort()
		return
	}

	c.logger.Info(c.context, "Handle get profile success, profile: ", string(p))
	context.JSON(http.StatusOK, gin.H{
		"code": errs.SUCCESS,
		"msg":  errs.GetMsg(errs.SUCCESS),
		"data": string(p),
	})
}

func (c *Client) UpdateProfile(context *gin.Context) {
	profile := &Profile{}
	if err := context.ShouldBind(profile); err != nil {
		c.logger.Error(c.context, "Bind request data error, err: ", err.Error())
		context.JSON(http.StatusBadRequest, gin.H{
			"code": errs.ERR_UPDATE_PROFILE_FAILED,
			"msg":  errs.GetMsg(errs.ERR_UPDATE_PROFILE_FAILED),
			"data": nil,
		})
		context.Abort()
		return
	}

	userId := c.getAuthedData(context, KEY_USER_ID)
	if userId == nil {
		return
	}

	r := &userinfo.UpdateProfileRequest{
		Profile: &userinfo.Profile{
			UserId:   userId.(uint64),
			Username: profile.Username,
			Birthday: profile.Birthday,
		},
		RequestId: GetRequestId(context),
	}
	_, err := c.userinfoClient.UpdateProfile(context, r)
	if err != nil {
		c.logger.Error(c.context, "Call rpc server failed, error: ", err)
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

	c.logger.Info(c.context, "Handle update profile success.")
	context.JSON(http.StatusOK, gin.H{
		"code": errs.SUCCESS,
		"msg":  errs.GetMsg(errs.SUCCESS),
		"data": nil,
	})
}

func (c *Client) CreateProfile(context *gin.Context) {
	profile := &Profile{}
	if err := context.ShouldBind(profile); err != nil {
		c.logger.Error(c.context, "Bind request data error, err: ", err.Error())
		context.JSON(http.StatusBadRequest, gin.H{
			"code": errs.ERR_CREATE_PROFILE_FAILED,
			"msg":  errs.GetMsg(errs.ERR_CREATE_PROFILE_FAILED),
			"data": nil,
		})
		context.Abort()
		return
	}

	userId := c.getAuthedData(context, KEY_USER_ID)
	email := c.getAuthedData(context, KEY_EMAIL)
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
		RequestId: GetRequestId(context),
	}
	_, err := c.userinfoClient.CreateProfile(context, r)
	if err != nil {
		c.logger.Error(c.context, "Call rpc server failed, error: ", err)
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

	c.logger.Info(c.context, "Handle create profile success.")
	context.JSON(http.StatusOK, gin.H{
		"code": errs.SUCCESS,
		"msg":  errs.GetMsg(errs.SUCCESS),
		"data": nil,
	})
}

func (c *Client) getAuthedData(context *gin.Context, key string) any {
	userId, ok := context.Get(key)
	if !ok {
		c.logger.Error(c.context, "Get cookie failed, key: ", key)
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
