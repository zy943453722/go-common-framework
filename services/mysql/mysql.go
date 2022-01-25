package mysql

import (
	"strconv"
	"strings"
	"sync"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/server/web"
	"go-common-framework/services/log"
)

var logMux sync.Mutex

func InitDB() error {
	user, err := web.AppConfig.String("mysql.user")
	if err != nil {
		return err
	}
	password, err := web.AppConfig.String("mysql.password")
	if err != nil {
		return err
	}
	masterServer, err := web.AppConfig.String("mysql.hostmaster")
	if err != nil {
		return err
	}
	slaveServer, err := web.AppConfig.String("mysql.hostslave")
	if err != nil {
		return err
	}
	port, err := web.AppConfig.Int("mysql.port")
	if err != nil {
		return err
	}
	dbname, err := web.AppConfig.String("mysql.dbname")
	if err != nil {
		return err
	}
	masterString := user + ":" + password + "@tcp(" + masterServer + ":" + strconv.Itoa(port) + ")/" + dbname + "?charset=utf8"
	slaveString := user + ":" + password + "@tcp(" + slaveServer + ":" + strconv.Itoa(port) + ")/" + dbname + "?charset=utf8"
	if err = orm.RegisterDataBase("write", "mysql", masterString); err != nil {
		return err
	}
	if err = orm.RegisterDataBase("read", "mysql", slaveString); err != nil {
		return err
	}
	if err = orm.RegisterDataBase("default", "mysql", slaveString); err != nil {
		return err
	}
	orm.SetMaxIdleConns("write", 30)
	orm.SetMaxIdleConns("read", 30)
	orm.SetMaxIdleConns("default", 30)
	orm.Debug = true
	return nil
}

func InitOrmLogFunc(requestId string) {
	logMux.Lock()
	defer logMux.Unlock()
	orm.LogFunc = func(query map[string]interface{}) {
		log.Info(requestId, "[%s] [cost:%vms] [sql:%s]", strings.Trim(query["flag"].(string), " "), query["cost_time"], query["sql"])
	}
}
