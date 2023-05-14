package controllers

import (
	"github.com/beego/beego/v2/core/logs"
	"strconv"
	"web/models"
)

// RelationController 关系控制器
type RelationController struct {
	BaseController
}
// Index 首页
func (c *RelationController) Index() {
	c.Prepare()
	c.TplName = "relation/server_list.html"


	webHookId,err := strconv.Atoi(c.Ctx.Input.Param(":id"))

	if err != nil || webHookId <= 0{
		c.ServerError("WebHook does not exist.")
	}
	webHook := models.NewWebHook()
	webHook.WebHookId = webHookId

	if err := webHook.Find(); err != nil {
		c.ServerError("WebHook does not exist." )
	}
	if webHook.CreateAt != c.User.UserId {
		c.Forbidden("")
	}

	c.Data["Model"] = webHook

	res,err := models.NewRelation().QueryByWebHookId(webHookId,c.User.UserId)

	c.Data["lists"] = res
	c.Data["WebHook"] = true
}

// AddServer 检索服务器并添加到数据库
func (c *RelationController) AddServer() {
	c.Prepare()

	if c.Ctx.Input.IsPost() {
		webHookId ,err := c.GetInt("web_hook_id")

		if err != nil {
			c.JsonResult(500,"Parameter error: web_hook_id is require.")
		}
		serverParams := c.GetStrings("server_id")
		if len(serverParams) <= 0 {
			c.JsonResult(500,"Server Id is require.")
		}

		webHook := models.NewWebHook()
		webHook.WebHookId = webHookId

		if err := webHook.Find();err != nil {
			c.JsonResult(404,"WebHook does not exist.")
		}
		if webHook.CreateAt != c.User.UserId {
			c.JsonResult(403,"Permission denied")
		}

		serverIds := make([]int,len(serverParams))
		index := 0

		for _,id := range serverParams {
			if id,err := strconv.Atoi(id);err == nil {
				serverIds[index] = id
				index++
			}

		}
		servers,err := models.NewServer().QueryServerByServerId(serverIds,c.User.UserId)

		if err != nil {
			c.JsonResult(500,"An error occurred while querying data")
		}

		if len(servers) <= 0 {
			c.JsonResult(500,"添加的服务器无效")
		}

		relations := make([]map[string]interface{},len(servers))

		index = 0
		for _,server := range servers {
			relation := models.NewRelation()

			relation.WebHookId = webHookId
			relation.ServerId = server.ServerId
			relation.MemberId = c.User.UserId

			if err := relation.Save();err == nil {
				temp := map[string]interface{} {
					"server_id" : server.ServerId,
					"name"		: server.Name,
					"type"		: server.Type,
					"ip_address"	: server.IpAddress,
					"port"		: server.Port,
					"add_time"	: relation.CreateTime,
					"relation_id"	: relation.WebHookId,
					"status"	: server.Status,

				}
				relations[index] = temp
				index++
			}
		}

		c.JsonResult(0,"ok",relations)
	}

	keyword := c.GetString("keyword","")

	if keyword == "" {
		c.JsonResult(500,"Keyword is require.")
	}

	webHookId ,_ := c.GetInt("id")

	relation := models.NewRelation()

	var serverIds []int

	if relations ,err := relation.QueryByWebHookId(webHookId,c.User.UserId);err == nil && len(relations) > 0{
		serverIds = make([]int,len(relations))
		i := 0
		for _,item := range relations {
			serverIds[i] = item.ServerId
			i++
		}
	}



	serverList,err := models.NewServer().Search(keyword,c.User.UserId,serverIds...)

	if err != nil {
		c.JsonResult(500,"Query Result Error")
	}

	c.JsonResult(0,"ok",serverList)

	c.StopRun()
}

// DeleteServer 删除一个服务器
func (c *RelationController) DeleteServer() {
	c.Prepare()
	relationId,err := strconv.Atoi(c.Ctx.Input.Param(":id"))

	if err != nil {
		c.JsonResult(500,"Parameter error: web_hook_id is require.")
	}

	relation := models.NewRelation()

	if err := relation.Find(relationId);err != nil {
		logs.Info("DeleteServer:",err)

		c.JsonResult(404,"Server does not exist.")
	}

	server := models.NewServer()
	server.ServerId = relation.ServerId

	if err := server.Find();err != nil || server.CreateAt != c.User.UserId {
		c.JsonResult(403,"Permission denied")
	}
	webHook := models.NewWebHook()
	webHook.WebHookId = relation.WebHookId

	if err := webHook.Find();err != nil || webHook.CreateAt != c.User.UserId{
		c.JsonResult(403,"Permission denied")
	}

	if err := relation.Delete() ; err != nil{
		c.JsonResult(500,"Delete failed")
	}

	models.NewScheduler().DeleteByWhere(" AND relation_id  = ?", relationId)

	c.JsonResult(0,"ok")
}
