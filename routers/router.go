package routers

import (
	"github.com/beego/beego/v2/server/web"
	"go-common-framework/controllers"
)

const (
	RUN_MODE_DEV  = "dev"
	RUN_MODE_TEST = "test"
	RUN_MODE_PRE  = "pre"
	RUN_MODE_PROD = "prod"
)

func init() {
	router := web.NewNamespace("/api",
		web.NSNamespace("/consumer",
			web.NSRouter("/list", &controllers.ConsumerController{}, "get:GetConsumerList"),
			web.NSRouter("/export", &controllers.ConsumerController{}, "get:ExportConsumerList"),
			web.NSRouter("/info", &controllers.ConsumerController{}, "get:GetConsumerInfo"),
			web.NSRouter("/edit", &controllers.ConsumerController{}, "put:EditConsumer"),
		),
	)
	web.AddNamespace(router)
	web.InsertFilter("/api/consumer/*", web.BeforeExec, LoginTestFilter)
}
