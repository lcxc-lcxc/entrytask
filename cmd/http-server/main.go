package main

import (
	"entrytask/global"
	"entrytask/internel/conn"
	"entrytask/internel/web"
	"entrytask/pkg/setting"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

var (
	config string
	port   string
	mode   string
)

func init() {
	cpuNum := runtime.NumCPU()
	runtime.GOMAXPROCS(cpuNum - 1)
	err := setupFlag()
	if err != nil {
		log.Fatalf("HTTP setup falied")
	}
	err = setupSetting()
	if err != nil {
		log.Fatalf("HTTP setup falied")
	}

	err = setupDBEngine()
	if err != nil {
		log.Fatalf("HTTP Set up DBEngine fail %v\n", err)
	}

	err = setupRPCClient()
	if err != nil {
		log.Fatalf("HTTP Set up RPC Client fail: %v\n", err)
	}

	err = setupCacheClient()
	if err != nil {
		log.Fatalf("HTTP Set up Cache Client fail: %v\n", err)
	}

	err = setupLogger()
	if err != nil {
		log.Fatalf("HTTP Set up Logger fail: %v\n", err)
	}

}

func setupDBEngine() error {
	var err error
	global.DBEngine, err = conn.NewDBEngine(global.DBSetting)
	if err != nil {
		log.Println("Set up DBEngine fail")
		return err
	}
	log.Println("Set up DBEngine Success")
	return nil

}

//const (
//	Ldate         = 1 << iota     // the date in the local time zone: 2009/01/23
//	Ltime                         // the time in the local time zone: 01:23:23
//	Lmicroseconds                 // microsecond resolution: 01:23:23.123123.  assumes Ltime.
//	Llongfile                     // full file name and line number: /a/b/c/d.go:23
//	Lshortfile                    // final file name element and line number: d.go:23. overrides Llongfile
//	LUTC                          // if Ldate or Ltime is set, use UTC rather than the local time zone
//	Lmsgprefix                    // move the "prefix" from the beginning of the line to before the message
//	LstdFlags     = Ldate | Ltime // initial values for the standard logger
//)
func setupLogger() error {
	logFile, err := os.OpenFile("log/httpserver.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("open log file error :%v\n", err)
		return err
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	return nil
}

func setupCacheClient() error {
	var err error
	global.RedisClient, err = conn.NewCacheClient(global.CacheSetting)
	if err != nil {
		log.Fatalf("Set up redis Client fail")
		return err
	}
	log.Println("Set up redis Client Success")
	return nil
}

func setupRPCClient() error {
	var err error
	global.GRPCClient, err = conn.NewRPCClient(global.RpcClientSetting)
	if err != nil {
		log.Fatalf("Set up RPC Client fail")
		return err
	}
	log.Println("Set up RPC Client Success")
	return nil
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
