package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/net/context"

	"gitlab.sftcwl.com/golang-lib-inner/golang/toml"
	"gitlab.sftcwl.com/tc-inf/unique-number/server"
)

func main() {
	/* ./sevice -config ./conf/config.toml启动服务*/
	config := flag.String("config", "./conf/config.toml", "config file")
	flag.Parse()
	//配置读取工具
	cfg := toml.NewTomlConfig()
	c := server.Config{}
	if err := cfg.Read(*config, &c); err != nil {
		fmt.Println("Read err:", err)
		os.Exit(0)
	}
	ctx, cancel := context.WithCancel(context.Background())
	//配置加载
	if err := cfg.InitConfig(ctx); err != nil {
		fmt.Println("InitConfig err:", err)
		os.Exit(0)
	}
	//初始化
	if err := server.InitServer(ctx, cfg, c); err != nil {
		fmt.Println("InitServer err:", err)
		os.Exit(0)
	}
	//监听退出信号
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	fmt.Println("Server done")
	cancel()
}
