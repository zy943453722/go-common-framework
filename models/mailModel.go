package models

import (
	"fmt"
	"time"

	"github.com/beego/beego/v2/server/web"
	"go-common-framework/services/log"
	"go-common-framework/services/mail"
)

var exceptionMailTemplate = `
		<p style="font-size:20px">Hi, all</p>
		<p style="font-size:20px">时间: %v</p>
		<p style="font-size:20px">环境: %s</p>
		<p style="font-size:20px">requestId: %s</p>
		<p style="font-size:20px">异常: %s</p>
	`

func HandleExceptionMail(e interface{}, requestId, msg string, t time.Time) {
	errMsg := ""
	switch t := e.(type) {
	case string:
		errMsg = e.(string)
		log.Error(requestId, e)
	case error:
		errMsg = t.Error()
		log.Error(requestId, t.Error())
	default:
		errMsg = fmt.Sprintf("未知错误:%v", e)
		log.Error(requestId, errMsg)
	}
	msg += fmt.Sprintf(", 原因:%s", errMsg)
	mailInfo, err := mail.NewMail(nil, mail.ERROR_MAIL_SUBJECT, "",
		nil,
		[]byte(fmt.Sprintf(exceptionMailTemplate, t, web.BConfig.RunMode, requestId, msg)))
	if err != nil {
		log.Error(requestId, err.Error())
	}
	mail.MailChan <- mailInfo
}
