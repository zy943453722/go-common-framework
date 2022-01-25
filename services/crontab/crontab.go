package crontab

import (
	"context"
	"fmt"
	"go-common-framework/types/request"
	"time"

	"github.com/beego/beego/v2/server/web"

	"github.com/beego/beego/v2/task"
	"go-common-framework/models"
	"go-common-framework/services/log"
	"go-common-framework/services/redis"
	"go-common-framework/util"
)

func initCrontabConfig() (map[string]string, error) {
	var err error
	crontabConfig := make(map[string]string)
	crontabConfig["consumer"], err = web.AppConfig.String("crontab.consumer")
	if err != nil {
		return nil, err
	}
	return crontabConfig, nil
}

func InitCrontab() error {
	crontabConfig, err := initCrontabConfig()
	if err != nil {
		return err
	}
	syncConsumerTask := task.NewTask("syncConsumerTask", crontabConfig["consumer"], SyncConsumers)
	task.AddTask("syncConsumerTask", syncConsumerTask)
	task.StartTask()
	return nil
}

func SyncConsumers(ctx context.Context) error {
	requestId := util.GenUuid()
	msg := ""
	defer func() {
		if e := recover(); e != nil {
			models.HandleExceptionMail(e, requestId, msg, time.Now())
		}
	}()
	res, err := redis.RedisInstance.RedisSetExNx(requestId, redis.CONSUMER_LOCK, redis.LOCK_EXPIRED, "1")
	if err != nil {
		msg = fmt.Sprintf("获取客户失败,加redis分布式锁失败, key:%s", redis.CONSUMER_LOCK)
		panic(err)
	}
	if res == "" {
		log.Info(requestId, "获取客户redis加分布式锁失败,请关注其他机器")
		return nil
	}
	defer redis.RedisInstance.RedisDel(requestId, redis.CONSUMER_LOCK)

	models.GetConsumerInfo(requestId, &request.GetConsumerInfoRequest{
		Cid: "xxxx",
	})

	return nil
}
