package handler

import (
	"context"
	"loggers"
	"protos/userinfo"
)

type Client struct {
	context        context.Context
	userinfoClient userinfo.UserinfoService
	logger         *logger.Logger
}

func NewClient(context context.Context, userinfoClient userinfo.UserinfoService, logger *logger.Logger) *Client {
	return &Client{
		context:        context,
		userinfoClient: userinfoClient,
		logger:         logger,
	}
}
