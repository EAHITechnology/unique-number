package server

import (
	"github.com/EAHITechnology/inf/golang/enet"
	"github.com/EAHITechnology/inf/golang/log"
	"github.com/EAHITechnology/inf/golang/toml"
	"github.com/EAHITechnology/inf/unique-number/handler"
	"github.com/EAHITechnology/inf/unique-number/logic"
	"github.com/EAHITechnology/inf/unique-number/proto/pb"
	"golang.org/x/net/context"
)

func InitServer(ctx context.Context, cfg *toml.TomlConfig, c Config) error {
	if err := initConfig(c); err != nil {
		log.Errorf("initConfig err:%s", err.Error())
		return err
	}
	if err := handler.InitHandler(); err != nil {
		log.Errorf("InitManager err:%s", err.Error())
		return err
	}
	if err := logic.InitManager(ctx); err != nil {
		log.Errorf("InitManager err:%s", err.Error())
		return err
	}
	if err := initRpcServer(ctx, cfg); err != nil {
		log.Errorf("initRpcServer err:%s", err.Error())
		return err
	}
	return nil
}

func initRpcServer(ctx context.Context, cfg *toml.TomlConfig) error {
	s, lis, err := enet.IniRpctLis()
	if err != nil {
		return err
	}
	pb.RegisterUnServiceServer(s, &handler.GetUniqueNumberServer{})
	enet.ReflectionRegister(s)
	go func() {
		err = s.Serve((*lis))
		if err != nil {
			log.Errorf("Serve err:%s", err.Error())
		}
	}()
	log.Info("grpc serve init ok")
	return nil
}
