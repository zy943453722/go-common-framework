package main

import (
	"fmt"
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/task"
	_ "github.com/go-sql-driver/mysql"
	_ "go-common-framework/dao/api/serviceApi"
	_ "go-common-framework/routers"
	"go-common-framework/services/crontab"
	"go-common-framework/services/log"
	"go-common-framework/services/mail"
	"go-common-framework/services/mysql"
	"go-common-framework/services/redis"
	"net/http"
	_ "net/http/pprof"
	"os"
)

func main() {
	if err := log.InitLogger(); err != nil {
		fmt.Printf("init log error: %s", err.Error())
		os.Exit(1)
	}
	if err := mysql.InitDB(); err != nil {
		fmt.Printf("init mysql error:%s", err.Error())
		os.Exit(1)
	}
	if err := redis.InitRedis(); err != nil {
		fmt.Printf("init redis error:%s", err.Error())
		os.Exit(1)
	}
	if err := mail.InitMailPool(); err != nil {
		fmt.Printf("init mail pool error:%s", err.Error())
		os.Exit(1)
	}
	monitorPort, err := web.AppConfig.String("monitorport")
	if err != nil {
		fmt.Printf("get monitorPort config error:%s", err.Error())
		os.Exit(1)
	}

	go http.ListenAndServe(monitorPort, nil)

	go mail.SendMail()
	crontab.InitCrontab()
	defer task.StopTask()

	web.Run()
}
