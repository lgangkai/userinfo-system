package handler

import (
	"context"
	errs "errs"
	"github.com/asim/go-micro/v3/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"loggers"
	"net/http"
	"protos/userinfo"
)

const (
	KEY_REQUEST_ID     = "request_id"
	KEY_ACCESS_TOKEN   = "access_token"
	KEY_USER_ID        = "user_id"
	KEY_EMAIL          = "email"
	COOKIE_EXPIRE_TIME = 3600 * 24
)

func (c *Client) Authenticate(context *gin.Context) {
	c.logger.Info(c.context, "Authenticate for api request.")
	token, err := context.Cookie(KEY_ACCESS_TOKEN)
	if err != nil {
		c.logger.Error(c.context, "Get access_token from cookie failed, err: ", err.Error())
		context.JSON(http.StatusUnauthorized, gin.H{
			"code": errs.ERR_AUTH_FAILED,
			"msg":  errs.GetMsg(errs.ERR_AUTH_FAILED),
			"data": nil,
		})
		context.Abort()
		return
	}

	req := &userinfo.AuthRequest{
		Token:     token,
		RequestId: GetRequestId(context),
	}
	resp, err := c.userinfoClient.Authenticate(context, req)
	if err != nil {
		c.logger.Error(c.context, "Authenticate failed, err: ", err.Error())
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

	c.logger.Info(c.context, "Authenticate succeed, userId: ", resp.GetUserId())
	context.Set(KEY_USER_ID, resp.GetUserId())
	context.Set(KEY_EMAIL, resp.GetEmail())
	context.Next()
}

func (c *Client) Log(context *gin.Context) {
	c.logger.Info(c.context, "Handling request: ", context.FullPath(), " method: ", context.Request.Method)
}

// GenRequestId assign an unique id to each request.
// Critical for live bug location.
func (c *Client) GenRequestId(ctx *gin.Context) {
	u, err := uuid.NewUUID()
	if err != nil {
		c.logger.Warning(c.context, "Gen request_id failed, err: ", err.Error())
		return
	}
	ctx.Set(KEY_REQUEST_ID, u.String())
}

func (c *Client) SetTraceData(ctx *gin.Context) {
	userId, ok := ctx.Get(KEY_USER_ID)
	if !ok {
		userId = uint64(0)
	}
	c.context = context.WithValue(ctx, logger.TraceDataKey{}, logger.TraceData{
		RequestId: GetRequestId(ctx),
		UserId:    userId.(uint64),
	})
}

func GetRequestId(ctx *gin.Context) string {
	requestId, ok := ctx.Get(KEY_REQUEST_ID)
	if !ok {
		requestId = ""
	}
	return requestId.(string)
}
