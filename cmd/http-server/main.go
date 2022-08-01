package main

import (
	"entrytask/global"
	"entrytask/internel/web"
	"entrytask/pkg/setting"
	"flag"
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	config string
	port   string
	mode   string
)

func init() {
	err := setupFlag()
	if err != nil {
		log.Fatalf("HTTP setup falied")
	}
	err = setupSetting()
	if err != nil {
		log.Fatalf("HTTP setup falied")
	}

}

func setupFlag() error {
	//StringVar defines a string flag with specified name, default value, and usage string. The argument p points to a string variable in which to store the value of the flag.
	flag.StringVar(&port, "port", "", "启动端口")
	flag.StringVar(&mode, "mode", "", "启动模式")
	flag.StringVar(&config, "config", "./config", "配置文件路径")
	flag.Parse()
	return nil
}

func setupSetting() error {
	log.Printf("%v", config)
	s, err := setting.NewSetting(strings.Split(config, ",")...)
	if err != nil {
		return err
	}
	err = s.ReadSection("RpcServer", &global.RpcServerSetting)
	if err != nil {
		return err
	}
	err = s.ReadSection("HttpServer", &global.HttpServerSetting)
	if err != nil {
		return err
	}
	err = s.ReadSection("Database", &global.DBSetting)
	if err != nil {
		return err
	}
	err = s.ReadSection("Redis", &global.CacheSetting)
	if err != nil {
		return err
	}
	err = s.ReadSection("RpcClient", &global.RpcClientSetting)
	if err != nil {
		return err
	}

	if port != "" {
		global.HttpServerSetting.Port = port
	}

	if mode != "" {
		global.HttpServerSetting.Mode = mode
	}
	return nil

}
func main() {
	r := web.NewRouter()
	server := &http.Server{
		Addr:              ":" + global.HttpServerSetting.Port,
		Handler:           r,
		ReadHeaderTimeout: global.HttpServerSetting.ReadTimeout * time.Second,
		WriteTimeout:      global.HttpServerSetting.WriteTimeout * time.Second,
	}

	log.Printf("Starting HTTP Server , Listening %v ", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Server ListenAndServe Fail %v", err)
	}

}
