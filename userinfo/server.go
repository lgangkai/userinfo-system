package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/asim/go-micro/plugins/registry/etcd/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/logger"
	"github.com/asim/go-micro/v3/registry"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/redis/go-redis/v9"
	"protos/userinfo"
	"user-server/conf"
	"user-server/dao"
	"user-server/wire"
)

type Server struct {
	service micro.Service
}

func (s *Server) Init() error {
	logger.Info("Init server...")
	// 1. load config file.
	var confPath string
	flag.StringVar(&confPath, "config", "conf/userinfo.yaml", "define config file")
	flag.Parse()
	config, err := conf.LoadConfig(confPath)
	if err != nil {
		logger.Error("load config file error, err: ", err)
		return err
	}
	logger.Info("config file loaded, config: ", config)

	mysqlMasterConf := config.MysqlMaster
	mysqlSlaveConf := config.MysqlSlave
	etcdConf := config.Etcd
	microConf := config.Micro
	redisConf := config.Redis

	// 2. register service.
	etcdReg := etcd.NewRegistry(
		registry.Addrs(etcdConf.Addrs...),
	)
	s.service = micro.NewService(
		micro.Name(microConf.Name),
		micro.Address(microConf.Addr),
		micro.Registry(etcdReg),
	)

	// 3. init basic dependencies.
	sqlMaster, err := sql.Open(mysqlMasterConf.Driver, fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", mysqlMasterConf.Name,
		mysqlMasterConf.Password, mysqlMasterConf.Host, mysqlMasterConf.Port, mysqlMasterConf.DB))
	if err != nil {
		logger.Error("init sqlDB master failed, err: ", err.Error())
	}

	sqlSlave, err := sql.Open(mysqlSlaveConf.Driver, fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", mysqlSlaveConf.Name,
		mysqlSlaveConf.Password, mysqlSlaveConf.Host, mysqlSlaveConf.Port, mysqlSlaveConf.DB))
	if err != nil {
		logger.Error("init sqlDB slave failed, err: ", err.Error())
	}

	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: redisConf.Addrs,
	})

	// 4. injection.
	userinfoHandler := wire.InitUserinfoHandler(
		&dao.DBMaster{DB: sqlMaster},
		&dao.DBSlave{DB: sqlSlave},
		rdb,
	)

	// 5. init service
	s.service.Init()
	err = userinfo.RegisterUserinfoHandler(s.service.Server(), userinfoHandler)
	if err != nil {
		logger.Error("register UserinfoHandler failed, err: ", err.Error())
		return err
	}

	return nil
}

func (s *Server) Run() error {
	if err := s.service.Run(); err != nil {
		logger.Error("run server failed, err: ", err.Error())
		return err
	}
	return nil
}
