package handler

import (
	errs "errs"
	"github.com/asim/go-micro/v3/errors"
	"github.com/asim/go-micro/v3/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"protos/userinfo"
)

const (
	KEY_ACCESS_TOKEN   = "access_token"
	KEY_USER_ID        = "user_id"
	KEY_EMAIL          = "email"
	COOKIE_EXPIRE_TIME = 3600 * 24
)

func (c *Client) Authenticate(context *gin.Context) {
	logger.Info("Authenticate for api request: ", context.FullPath())
	token, err := context.Cookie(KEY_ACCESS_TOKEN)
	if err != nil {
		logger.Error("Get access_token from cookie failed, err: ", err.Error())
		context.JSON(http.StatusUnauthorized, gin.H{
			"code": errs.ERR_AUTH_FAILED,
			"msg":  errs.GetMsg(errs.ERR_AUTH_FAILED),
			"data": nil,
		})
		context.Abort()
		return
	}

	req := &userinfo.AuthRequest{Token: token}
	resp, err := c.UserinfoClient.Authenticate(context, req)
	if err != nil {
		logger.Error("Authenticate failed, err: ", err.Error())
		code := errors.Parse(err.Error()).Code
		msg := errors.Parse(err.Error()).Detail
		context.JSON(http.StatusUnauthorized, gin.H{
			"code": code,
			"msg":  msg,
			"data": nil,
		})
		context.Abort()
		return
	}

	logger.Info("Authenticate succeed, userId: ", resp.GetUserId())
	context.Set(KEY_USER_ID, resp.GetUserId())
	context.Set(KEY_EMAIL, resp.GetEmail())
	context.Next()
}

func Log(context *gin.Context) {
	logger.Info("Handling request: ", context.FullPath(), " method: ", context.Request.Method)
}
