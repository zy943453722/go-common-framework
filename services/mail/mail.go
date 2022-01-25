package mail

import (
	"net/smtp"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/server/web"
	"github.com/jordan-wright/email"
	"go-common-framework/services/log"
	"go-common-framework/util"
)

const (
	MAIL_CHAN_LENGTH  = 10
	MAIL_CLIENT       = 2
	MAIL_CONN_TIMEOUT = 10

	ERROR_MAIL_SUBJECT = "xxxxxx"
)

var (
	//邮件池
	MailChan = make(chan *email.Email, MAIL_CHAN_LENGTH)
	//邮件客户端池
	MailClientPool = new(email.Pool)
)

func InitMailPool() error {
	host, err := web.AppConfig.String("mail.host")
	if err != nil {
		return err
	}
	port, err := web.AppConfig.Int("mail.port")
	if err != nil {
		return err
	}
	username, err := web.AppConfig.String("mail.username")
	if err != nil {
		return err
	}
	password, err := web.AppConfig.String("mail.password")
	if err != nil {
		return err
	}
	/*MailClientPool, err = email.NewPool(host+":"+strconv.Itoa(port), MAIL_CLIENT, LoginAuth(username, password), &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return err
	}*/
	MailClientPool, err = email.NewPool(host+":"+strconv.Itoa(port), MAIL_CLIENT, smtp.PlainAuth("", username, password, host))
	if err != nil {
		return err
	}

	return nil
}

func NewMail(cc []string, subject, fileName string, textContent, htmlContent []byte) (*email.Email, error) {
	mailInfo := email.NewEmail()
	from, err := web.AppConfig.String("mail.from")
	if err != nil {
		return nil, err
	}
	mailInfo.From = from
	to, err := web.AppConfig.String("mail.to")
	if err != nil {
		return nil, err
	}
	mailInfo.To = strings.Split(to, ",")
	if cc != nil {
		mailInfo.Cc = cc
	}
	mailInfo.Subject = subject
	if fileName != "" {
		if _, err = mailInfo.AttachFile(fileName); err != nil {
			return nil, err
		}
	}
	if textContent != nil {
		mailInfo.Text = textContent
	}
	if htmlContent != nil {
		mailInfo.HTML = htmlContent
	}
	return mailInfo, nil
}

func SendMail() {
	requestId := util.GenUuid()
	for mailInfo := range MailChan {
		go func(mailInfo *email.Email) {
			if err := MailClientPool.Send(mailInfo, MAIL_CONN_TIMEOUT*time.Second); err != nil {
				log.Error(requestId, err.Error())
			}
		}(mailInfo)
	}
}
