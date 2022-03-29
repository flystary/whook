package cli

import (
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"net/url"
	"time"
	"web/models"
	"web/tasks"
)

// RegisterDataBase 注册数据库
func RegisterDataBase() {
	host, _ := beego.AppConfig.String("db_host")
	database,_ := beego.AppConfig.String("db_database")
	username,_ := beego.AppConfig.String("db_username")
	password,_ := beego.AppConfig.String("db_password")
	timezone,_ :=  beego.AppConfig.String("timezone")

	port,_ := beego.AppConfig.String("db_port")

	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true&loc=%s",username,password,host,port,database,url.QueryEscape(timezone))
	orm.RegisterDataBase("default", "mysql",dataSource)
	orm.DefaultTimeLoc, _ = time.LoadLocation(timezone)
}

// RegisterModel 注册Model
func RegisterModel()  {
	orm.RegisterModel(new(models.User))
	orm.RegisterModel(new(models.Server))
	orm.RegisterModel(new(models.WebHook))
	orm.RegisterModel(new(models.Scheduler))
	orm.RegisterModel(new(models.Relation))
}


// RegisterLogger 注册日志
func RegisterLogger()  {
	logs.SetLogger("console")
	logs.SetLogger(logs.AdapterConsole)
	logs.SetLogger(logs.AdapterFile,`{"filename":"logs/web.log","level":7,"maxlines":1000,"maxsize":1000,"daily":true,"maxdays":10,"color":true}`)
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(2)
//logger.Async()
//1-7级别递减，默认是trace，显示当前数字以前的级别，例如：3时，显示【Emergency】【Alert】【Critical】【Error】
//logger.Emergency("log3--->Emergency")
//logs.Alert("log3--->Alert")       //1
//logs.Critical("log3--->Critical") //2
//logs.Error("log3--->Error")       //3
//logs.Warn("log3--->Warning")      //4
//logs.Notice("log3--->Notice")     //5
//logs.Info("log3--->Info")         //6
//logs.Debug("log3--->Debug")       //7
//logger.Trace("log3--->Trace")
}

//// RegisterJsonLogger json日志
//func RegisterJsonLogger()  {
//	var logCfg = log.LoggerConfig{
//		FileName:            "logs/web_json.log",
//		Level:               7,
//		EnableFuncCallDepth: true,
//		LogFuncCallDepth:    3,
//		RotatePerm:          "777",
//		Perm:                "777",
//		Color:               false,
//	}
//
//	// 设置beego log库的配置
//	b, _ := json.Marshal(&logCfg)
//	logs.SetLogger(logs.AdapterConsole)
//	logs.SetLogger(logs.AdapterFile, string(b))
//	logs.EnableFuncCallDepth(logCfg.EnableFuncCallDepth)
//	logs.SetLogFuncCallDepth(logCfg.LogFuncCallDepth)
//
//	logs.Info("program start...")
//	logs.Critical("log3--->Critical") //2
//	logs.Error("log3--->Error")       //3
//	logs.Warn("log3--->Warning")      //4
//	logs.Notice("log3--->Notice")     //5
//	logs.Info("log3--->Info")         //6
//	logs.Debug("log3--->Debug")
//}

// RegisterTaskQueue 注册队列
func RegisterTaskQueue()  {

	schedulerList,err := models.NewScheduler().QuerySchedulerByState("wait");
	if err == nil {
		for _,scheduler := range schedulerList {
			tasks.Add(tasks.Task{SchedulerId: scheduler.SchedulerId,WebHookId:scheduler.WebHookId,ServerId:scheduler.ServerId})
		}
	}else{
		fmt.Println(err)
	}

}

//RunCli 注册orm命令行工具
func RunCli()  {
	orm.RunCommand()
	Install()
	Version()
}

//Run 启动web
func Run() {
	beego.Run()
}