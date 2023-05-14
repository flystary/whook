package controllers

import (
	"bytes"
	"github.com/beego/beego/v2/core/logs"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
	"time"
	"web/models"
	"web/modules/pager"
	"web/tasks"
)

// upgrader 默认的 WebSocket 选项
var upgrader = websocket.Upgrader{CheckOrigin : verification} // use default options

// SchedulerController 任务调度控制器
type SchedulerController struct {
	BaseController
}

// Index 首页
func (c *SchedulerController) Index()  {
	c.Prepare()
	c.TplName = "scheduler/index_vue.html"


	pageIndex, _ := c.GetInt("page", 1)
	relationId,_ := strconv.Atoi(c.Ctx.Input.Param(":id"))

	if relationId <= 0 {
		c.Abort("404")
	}

	relationDetailedResult ,err := models.FindRelationDetailedByWhere("AND relation_id = ?", relationId)

	if err != nil {
		logs.Info("FindRelationDetailed Error : ", err.Error())
		c.Abort("500")
	}
	var relationDetailed models.RelationDetailed

	if len(relationDetailedResult) > 0{
		relationDetailed = relationDetailedResult[0]
	}

	var schedulers []models.Scheduler

	pageOptions := pager.PageOptions{
		TableName:  models.NewScheduler().TableName(),
		EnableFirstLastLink : true,
		CurrentPage : pageIndex,
		PageSize : 15,
		ParamName : "page",
		Conditions : " AND relation_id = " + strconv.Itoa(relationDetailed.RelationId) + " order by scheduler_id desc",
	}

	totalItem, totalCount, rs, pageHtml := pager.GetPagerLinks(&pageOptions, c.Ctx)

	_,err = rs.QueryRows(&schedulers)      //把当前页面的数据序列化进一个切片内

	if err != nil {
		logs.Error("",err.Error())
	}

	var webList []models.WebScheduler

	if c.IsAjax() {
		if len(schedulers) > 0 {
			webList = make([]models.WebScheduler,len(schedulers))
			for i, item := range schedulers {
				webList[i] = (&item).ToWebScheduler()
			}
		}
		c.JsonResult(0,"ok",webList)
	}


	c.Data["Model"] = relationDetailed

	c.Data["lists"] = webList
	c.Data["html"] = pageHtml
	c.Data["totalItem"] = totalItem
	c.Data["totalCount"] = totalCount
	c.Data["WebHook"] = true
	c.Data["WebSocketUrl"] = "ws://" + c.Ctx.Request.Host + c.URLFor("SchedulerController.Status",":scheduler_id","")

}

// Console 控制台
func (c *SchedulerController) Console() {
	schedulerId,err := strconv.Atoi(c.Ctx.Input.Param(":scheduler_id"))

	if err != nil {
		c.JsonResult(500,"Parameter error")
	}

	scheduler := models.NewScheduler()
	scheduler.SchedulerId = schedulerId

	if err := scheduler.Find();err != nil {
		c.JsonResult(500,"Error 50001: Query data error")
	}
	deailed,err := models.FindRelationDetailedByWhere("AND relation_id = ? AND user_id = ?", scheduler.RelationId,c.User.UserId)

	if err != nil || len(deailed) <= 0{
		c.JsonResult(404,"The data does not exist")
	}

	data := map[string]interface{}{"log": scheduler.LogContent,"status": scheduler.Status}

	c.JsonResult(0,"ok",data)
}

// Cancel 取消任务
func (c *SchedulerController) Cancel() {
	schedulerId,err := strconv.Atoi(c.Ctx.Input.Param(":scheduler_id"))

	if err != nil {
		c.JsonResult(500,"Parameter error")
	}

	scheduler := models.NewScheduler()
	scheduler.SchedulerId = schedulerId

	if err := scheduler.Find();err != nil {
		c.JsonResult(500,"Error 50001: Query data error")
	}
	deailed,err := models.FindRelationDetailedByWhere("AND relation_id = ? AND user_id = ?", scheduler.RelationId,c.User.UserId)

	if err != nil || len(deailed) <= 0{
		c.JsonResult(404,"The data does not exist")
	}

	scheduler.Status = "suspend"

	if err := scheduler.Save(); err != nil {
		c.JsonResult(500,"Cancel failed")
	}
	c.JsonResult(0,"ok")
}

// Resume 重新执行
func (c *SchedulerController) Resume () {
	schedulerId,err := strconv.Atoi(c.Ctx.Input.Param(":scheduler_id"))

	if err != nil {
		c.JsonResult(500,"Parameter error")
	}

	scheduler := models.NewScheduler()
	scheduler.SchedulerId = schedulerId

	if err := scheduler.Find();err != nil {
		c.JsonResult(500,"Error 50001: Query data error")
	}
	deailed,err := models.FindRelationDetailedByWhere("AND relation_id = ? AND user_id = ?", scheduler.RelationId,c.User.UserId)

	if err != nil || len(deailed) <= 0{
		c.JsonResult(404,"The data does not exist")
	}
	newScheduler := models.NewScheduler()

	newScheduler.Status = "wait"
	newScheduler.ExecuteType = 1
	newScheduler.WebHookId = scheduler.WebHookId
	newScheduler.ServerId = scheduler.ServerId
	newScheduler.RelationId = scheduler.RelationId
	newScheduler.Data = scheduler.Data
	newScheduler.PushUser = scheduler.PushUser
	newScheduler.ShaValue	= scheduler.ShaValue


	if err := newScheduler.Save(); err != nil {
		c.JsonResult(500,"Cancel failed")
	}

	webModel := newScheduler.ToWebScheduler()


	go tasks.Add(tasks.Task{ SchedulerId : newScheduler.SchedulerId ,ServerId:newScheduler.ServerId,WebHookId:newScheduler.WebHookId})

	view ,_:= c.ExecuteViewPathTemplate("scheduler/index_item.html",webModel)

	data := map[string]interface{}{"view" : view,"data": webModel}
	c.JsonResult(0,"ok",data)
}

// Status 任务状态
func (c *SchedulerController) Status() {

	c.TplName = "scheduler/status.html"

	ws, err := upgrader.Upgrade(c.Ctx.ResponseWriter.ResponseWriter, c.Ctx.Request, nil)

	if err != nil {

		c.StopRun()
	}
	defer ws.Close()

	buf := bytes.NewBufferString("")

	schedulerId, err := strconv.Atoi(c.Ctx.Input.Param(":scheduler_id"))

	if err != nil {
		buf.WriteString("Parameter error.")
		ws.WriteMessage(websocket.TextMessage,buf.Bytes())
		c.StopRun()
	}
	scheduler := models.NewScheduler()

	scheduler.SchedulerId = schedulerId

	if err := scheduler.Find();err != nil {
		ws.WriteMessage(websocket.TextMessage,[]byte("Error 50001: Query data error"))
		c.StopRun()
	}
	detailed,err := models.FindRelationDetailedByWhere("AND relation_id = ? AND user_id = ?", scheduler.RelationId,c.User.UserId)

	if err != nil || len(detailed) <= 0{
		ws.WriteMessage(websocket.TextMessage,[]byte("The data does not exist"))
		c.StopRun()
	}
	for{
		m := scheduler.ToWebScheduler()

		err := ws.WriteJSON(m)
		if err != nil {
			c.StopRun()
		}

		time.Sleep(time.Second*2)

		scheduler = models.NewScheduler()

		scheduler.SchedulerId = schedulerId
		scheduler.Find()
	}
	c.StopRun()
}

// verification 校验是否可以连接 WebSocket
func verification( *http.Request) bool {

	return true
}
