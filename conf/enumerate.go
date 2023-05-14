package conf

import (
	beego "github.com/beego/beego/v2/server/web"
)

var LoginSessionName = "LoginSessionName"

func GetAppKey() string {
	return beego.AppConfig.DefaultString("app_key","web")
}

func QueueSize() int {
	queneSize := beego.AppConfig.DefaultInt("queue_size", 100)
	return queneSize
}
