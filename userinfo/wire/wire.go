//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"loggers"
	"user-server/biz/account"
	"user-server/biz/profile"
	"user-server/dao"
	"user-server/handler"
	account2 "user-server/service/account"
	profile2 "user-server/service/profile"
)

func InitUserinfoHandler(*dao.DBMaster, *dao.DBSlave, *redis.ClusterClient, *logger.Logger) *handler.UserinfoHandlerImpl {
	wire.Build(dao.NewProfileDao, profile.NewProfileBiz, profile2.NewProfileService, account.NewAccountBiz, account2.NewAccountService, dao.NewUserDao, handler.NewUserinfoHandlerImpl)
	return &handler.UserinfoHandlerImpl{}
}
