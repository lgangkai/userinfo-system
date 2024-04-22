package handler

import "protos/userinfo"

type Client struct {
	UserinfoClient userinfo.UserinfoService
}
