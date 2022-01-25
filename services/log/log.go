package log

import (
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
)

func InitLogger() error {
	level, err := web.AppConfig.Int("log.level")
	if err != nil {
		return err
	}

	filePath, err := web.AppConfig.String("log.file")
	if err != nil {
		return err
	}

	maxDays, err := web.AppConfig.Int("log.days")
	if err != nil {
		return err
	}

	logConfigStr := fmt.Sprintf("{\"filename\": \"%s\", \"level\":%d, \"maxdays\":%d }",
		filePath, level, maxDays)
	if err = logs.SetLogger(logs.AdapterFile, logConfigStr); err != nil {
		return err
	}
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(4)

	return nil
}

func formatString(args []interface{}) string {
	if len(args) > 1 {
		return fmt.Sprintf(args[0].(string), args[1:]...)
	}
	return args[0].(string)
}

func Info(requestId string, args ...interface{}) {
	str := formatString(args)
	str = fmt.Sprintf("[requestId = %s] ", requestId) + str
	logs.Info(str)
}

func Error(requestId string, args ...interface{}) {
	str := formatString(args)
	str = fmt.Sprintf("[requestId = %s] ", requestId) + str
	logs.Error(str)
}

func Warn(requestId string, args ...interface{}) {
	str := formatString(args)
	str = fmt.Sprintf("[requestId = %s] ", requestId) + str
	logs.Warn(str)
}

func Alert(requestId string, args ...interface{}) {
	str := formatString(args)
	str = fmt.Sprintf("[requestId = %s] ", requestId) + str
	logs.Alert(str)
}

func Critical(requestId string, args ...interface{}) {
	str := formatString(args)
	str = fmt.Sprintf("[requestId = %s] ", requestId) + str
	logs.Critical(str)
}

func Debug(requestId string, args ...interface{}) {
	str := formatString(args)
	str = fmt.Sprintf("[requestId = %s] ", requestId) + str
	logs.Debug(str)
}

func Trace(requestId string, args ...interface{}) {
	str := formatString(args)
	str = fmt.Sprintf("[requestId = %s] ", requestId) + str
	logs.Trace(str)
}

func Notice(requestId string, args ...interface{}) {
	str := formatString(args)
	str = fmt.Sprintf("[requestId = %s] ", requestId) + str
	logs.Notice(str)
}
