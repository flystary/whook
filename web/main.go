package main

import (
	_ "github.com/go-sql-driver/mysql"
	"web/cli"
	_ "web/modules/filters"
	_ "web/routers"
)

func main()	{
	//默认log
	cli.RegisterLogger()
	//自定义log
	//cli.RegisterJsonLogger()
	cli.RegisterDataBase()
	cli.RegisterModel()
	cli.RunCli()

	cli.RegisterTaskQueue()

	cli.Run()
}

