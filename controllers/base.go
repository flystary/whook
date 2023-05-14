package controllers

import (
	"bytes"
	beego "github.com/beego/beego/v2/server/web"
	"web/conf"
	"web/models"
)

type BaseController struct {
	beego.Controller
	User 	*models.User
	Scheme 	string
}

//func (c *BaseController) Get() {
//	c.Data["Website"] = "beego.me"
//	c.Data["Email"] = "astaxie@gmail.com"
//	c.TplName = "home/index.html"
//}

// Prepare 预处理.
func (c *BaseController) Prepare() {

	c.Data["SiteName"] = "web"
	c.Data["User"] = models.NewUser()


	if member,ok := c.GetSession(conf.LoginSessionName).(models.User); ok && member.UserId > 0{
		c.User = &member

		c.Data["User"] = c.User
	}


	scheme := "http"

	if c.Ctx.Request.TLS != nil {
		scheme = "https"
	}
	c.Scheme = scheme
}

// SetUser  获取或设置当前登录用户信息,如果 UserId 小于 0 则标识删除 Session
func (c *BaseController) SetUser(user models.User) {

	if user.UserId <= 0 {
		c.DelSession(conf.LoginSessionName)
		c.DelSession("uid")
		c.DestroySession()
	} else {
		c.SetSession(conf.LoginSessionName, user)
		c.SetSession("uid", user.UserId)
	}
}

// JsonResult 响应 json 结果
func (c *BaseController) JsonResult(errCode int,errMsg string, data ...interface{})  {
	json := make(map[string]interface{}, 3)

	json["errcode"] = errCode
	json["message"] = errMsg

	if len(data) > 0 && data[0] != nil {
		json["data"] = data[0]
	}

	c.Data["json"] = json
	c.ServeJSON(true)
	c.StopRun()
}

// UrlFor .
func (c *BaseController) UrlFor(endpoint string, values ...interface{}) string	 {
	return c.BaseUrl() + beego.URLFor(endpoint, values)

}

// BaseUrl .
func (c *BaseController) BaseUrl() string {
	if baseUrl, _ := beego.AppConfig.String("base_url"); baseUrl != "" {
		return baseUrl
	}
	return c.Ctx.Input.Scheme() + "://" + c.Ctx.Request.Host
}

// NotFound .
func (c *BaseController) NotFound(message interface{})  {
	c.TplName = "errors/404.html"
	c.Layout = ""
	c.Data["Model"] = map[string]interface{}{"Message":message}

	html,_ := c.RenderString()

	c.Abort(html)
}

// Forbidden .
func (c *BaseController) Forbidden(message interface{}) {
	c.TplName = "errors/403.html"
	c.Layout = ""
	c.Data["Model"] = map[string]interface{}{"Message":message}

	html,_ := c.RenderString()

	c.Abort(html)
}

// ServerError .
func (c *BaseController) ServerError (message interface{}) {
	c.TplName = "errors/500.html"
	c.Layout = ""
	c.Data["Model"] = map[string]interface{}{"Message":message}


	html,_ := c.RenderString()

	c.Abort(html)
}


// ExecuteViewPathTemplate 执行指定的模板并返回执行结果.
func (c *BaseController) ExecuteViewPathTemplate(tplName string,data interface{}) (string,error){
	var buf bytes.Buffer

	viewPath := c.ViewPath

	if c.ViewPath == "" {
		viewPath = beego.BConfig.WebConfig.ViewsPath

	}

	if err := beego.ExecuteViewPathTemplate(&buf,tplName,viewPath,data); err != nil {
		return "",err
	}
	return buf.String(),nil
}