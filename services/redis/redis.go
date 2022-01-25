package redis

import (
	"strconv"
	"time"

	_ "github.com/beego/beego/v2/client/cache/redis"
	"github.com/beego/beego/v2/server/web"
	"github.com/gomodule/redigo/redis"
	"go-common-framework/services/log"
)

const (
	CONSUMER_LOCK = "consumer_lock"

	LOCK_EXPIRED = 7200
)

var RedisInstance *RedisService

type RedisService struct {
	Pool *redis.Pool
}

func InitRedis() error {
	RedisInstance = new(RedisService)
	if err := RedisInstance.NewPool(); err != nil {
		return err
	}
	return nil
}

func (rs *RedisService) NewPool() error {
	host, err := web.AppConfig.String("redis.host")
	if err != nil {
		return err
	}
	port, err := web.AppConfig.Int("redis.port")
	if err != nil {
		return err
	}
	password, err := web.AppConfig.String("redis.password")
	if err != nil {
		return err
	}
	idleCount, err := web.AppConfig.Int("redis.idleCount")
	if err != nil {
		return err
	}
	activeCount, err := web.AppConfig.Int("redis.activeCount")
	if err != nil {
		return err
	}
	idleTimeout, err := web.AppConfig.Int("redis.idleTimeout")
	if err != nil {
		return err
	}
	connTimeoutConf, err := web.AppConfig.Int("redis.connTimeout")
	if err != nil {
		return err
	}
	redisPool := &redis.Pool{
		MaxIdle:     idleCount,
		MaxActive:   activeCount,
		IdleTimeout: time.Duration(idleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			redisAddress := host + ":" + strconv.Itoa(port)
			passwdOption := redis.DialPassword(password)
			connTimeout := redis.DialConnectTimeout(time.Duration(connTimeoutConf) * time.Millisecond)

			cc, err := redis.Dial("tcp", redisAddress, passwdOption, connTimeout)
			if err != nil {
				return nil, err
			}

			return cc, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < 10*time.Second {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}

	rs.Pool = redisPool
	return nil
}

func (rs *RedisService) RedisSet(requestId, key string, value string) (string, error) {
	redisConn := rs.Pool.Get()
	defer redisConn.Close()
	res, err := redis.String(redisConn.Do("SET", key, value))
	if err != nil {
		log.Error(requestId, "redis set failed, key[%s] value[%s] err[%s]", key, value, err)
	}
	return res, err
}

func (rs *RedisService) RedisSetEx(requestId, key string, expire int64, value string) (string, error) {
	redisConn := rs.Pool.Get()
	defer redisConn.Close()
	res, err := redis.String(redisConn.Do("SETEX", key, expire, value))
	if err != nil {
		log.Error(requestId, "redis setex failed, key[%s] value[%s] expire[%d] err[%s]", key, value, expire, err)
	}
	return res, err
}

func (rs *RedisService) RedisSetExNx(requestId, key string, expire int64, value string) (string, error) {
	redisConn := rs.Pool.Get()
	defer redisConn.Close()
	res, err := redis.String(redisConn.Do("SET", key, value, "NX", "EX", expire))
	if err == redis.ErrNil {
		return "", nil
	}
	if err != nil {
		log.Error(requestId, "redis setexnx failed, key[%s] value[%s] expire[%d] err[%s]", key, value, expire, err)
		return res, err
	}
	log.Info(requestId, "redis setexnx success, key[%s] res[%s]", key, res)
	return res, nil
}

func (rs *RedisService) RedisGet(requestId, key string) (string, error) {
	redisConn := rs.Pool.Get()
	defer redisConn.Close()
	res, err := redis.String(redisConn.Do("GET", key))
	if err != nil && err != redis.ErrNil {
		log.Error(requestId, "redis get failed, key[%s] err[%s]", key, err)
		return res, err
	}

	return res, nil
}

func (rs *RedisService) RedisDel(requestId, key string) (int64, error) {
	redisConn := rs.Pool.Get()
	defer redisConn.Close()
	res, err := redis.Int64(redisConn.Do("DEL", key))
	if err != nil {
		log.Error(requestId, "redis delete failed, key[%s] err[%s]", key, err)
	}

	log.Info(requestId, "redis delete success, key[%s]", key)
	return res, err
}

func (rs *RedisService) RedisIncr(requestId, key string) error {
	redisConn := rs.Pool.Get()
	defer redisConn.Close()
	_, err := redis.Bool(redisConn.Do("INCRBY", key, 1))
	if err != nil {
		log.Error(requestId, "redis incr failed, key[%s] err[%s]", key, err)
	}
	return err
}
