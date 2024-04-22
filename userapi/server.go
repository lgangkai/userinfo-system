package main

import (
	"flag"
	"github.com/asim/go-micro/plugins/registry/etcd/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/logger"
	"github.com/asim/go-micro/v3/registry"
	"protos/userinfo"
	"user-api/conf"
)

type Server struct {
	Addr           string
	UserinfoClient userinfo.UserinfoService
}

func (s *Server) Init() error {
	logger.Info("Init server...")
	// 1. load config file.
	var confPath string
	flag.StringVar(&confPath, "config", "conf/userapi.yaml", "define config file")
	flag.Parse()
	config, err := conf.LoadConfig(confPath)
	if err != nil {
		logger.Error("load config file error, err: ", err)
		return err
	}
	logger.Info("config file loaded, config: ", config)

	serverConf := config.Server
	microConf := config.Micro
	etcdConf := config.Etcd

	s.Addr = serverConf.Addr

	// 2. init microservice client
	etcdReg := etcd.NewRegistry(
		registry.Addrs(etcdConf.Addrs...),
	)
	mService := micro.NewService(micro.Registry(etcdReg))
	mService.Init()
	s.UserinfoClient = userinfo.NewUserinfoService(microConf.Name, mService.Client())

	return nil
}
