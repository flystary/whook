package controllers

import (
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"time"
	"web/conf"
	"web/models"
	"web/modules/gob"
)

// AccountController 用户登录与注册.
type AccountController struct {
		BaseController
}


// Login 用户登录.
func (c *AccountController) Login()  {
	c.Prepare()

	var reuser struct { UserId int ; Account string; Time time.Time}

	//如果Cookie中存在登录信息
	if cookie,ok := c.GetSecureCookie(conf.GetAppKey(),"login");ok{

		if err := gob.Decode(cookie,&reuser); err == nil {
			user := models.NewUser()
			user.UserId = reuser.UserId

			if err := models.NewUser().Find(); err == nil {
				c.SetUser(*user)

				c.Redirect(beego.URLFor("HomeController.Index"), 302)
				c.StopRun()
			}
		}
	}

	if c.Ctx.Input.IsPost() {
		account := c.GetString("inputAccount")
		password := c.GetString("inputPassword")

		users,err := models.NewUser().Login(account,password)

		//如果没有数据
		if err == nil {
			c.SetUser(*users)
			c.JsonResult(0,"ok")
			c.StopRun()
		}else{
			fmt.Println(err)
			c.JsonResult(500,"账号或密码错误",nil)
		}

		return
	}else{

		c.Layout = ""
		c.TplName = "account/login.html"
	}
}

// Logout 退出登录.
func (c *AccountController) Logout(){
	c.SetUser(models.User{})

	c.Redirect(beego.URLFor("AccountController.Login"),302)
}
