package handler

import (
	errs "errs"
	"github.com/asim/go-micro/v3/errors"
	"github.com/asim/go-micro/v3/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"protos/userinfo"
)

func (c *Client) Login(context *gin.Context) {
	account := &Account{}
	if err := context.ShouldBind(account); err != nil {
		logger.Error("Bind request data error, err: ", err.Error())
		context.JSON(http.StatusBadRequest, gin.H{
			"code": errs.ERR_LOGIN_REQUEST,
			"msg":  errs.GetMsg(errs.ERR_LOGIN_REQUEST),
			"data": nil,
		})
		context.Abort()
		return
	}
	r := &userinfo.LoginRequest{
		Email:    account.Email,
		Password: account.Password,
	}
	resp, err := c.UserinfoClient.Login(context, r)
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

	logger.Info("Handle login success.")
	context.SetCookie(KEY_ACCESS_TOKEN, resp.GetToken(), COOKIE_EXPIRE_TIME, "/", "", false, true)
	context.JSON(http.StatusOK, gin.H{
		"code": errs.SUCCESS,
		"msg":  errs.GetMsg(errs.SUCCESS),
		"data": nil,
	})
}

// Logout currently just delete the token and no need to call service.
func (c *Client) Logout(context *gin.Context) {
	context.SetCookie(KEY_ACCESS_TOKEN, "", -1, "/", "", false, true)
	context.JSON(http.StatusOK, gin.H{
		"code": errs.SUCCESS,
		"msg":  errs.GetMsg(errs.SUCCESS),
		"data": nil,
	})
}

func (c *Client) Register(context *gin.Context) {
	account := &Account{}
	if err := context.ShouldBind(account); err != nil {
		logger.Error("Bind request data error, err: ", err.Error())
		context.JSON(http.StatusBadRequest, gin.H{
			"code": errs.ERR_REGISTER_REQUEST,
			"msg":  errs.GetMsg(errs.ERR_REGISTER_REQUEST),
			"data": nil,
		})
		context.Abort()
		return
	}
	r := &userinfo.RegisterRequest{
		Email:    account.Email,
		Password: account.Password,
	}
	_, err := c.UserinfoClient.Register(context, r)
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

	logger.Info("Handle register success.")
	context.JSON(http.StatusOK, gin.H{
		"code": errs.SUCCESS,
		"msg":  errs.GetMsg(errs.SUCCESS),
		"data": nil,
	})
}
