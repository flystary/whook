package filters

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
)

func init() {
	Authorize := func(ctx *context.Context) {
		_, ok := ctx.Input.Session("uid").(int)
		if !ok {

			ctx.Redirect(302, "/login")
		}
	}
	beego.InsertFilter("/member/*",beego.BeforeRouter,Authorize);
	beego.InsertFilter("/member",beego.BeforeRouter,Authorize);
	beego.InsertFilter("/",beego.BeforeRouter,Authorize);
	beego.InsertFilter("/server",beego.BeforeRouter,Authorize)
	beego.InsertFilter("/server/*",beego.BeforeRouter,Authorize)
	beego.InsertFilter("/hook",beego.BeforeRouter,Authorize)
	beego.InsertFilter("/hook/*",beego.BeforeRouter,Authorize)
	beego.InsertFilter("/my",beego.BeforeRouter,Authorize)
}
