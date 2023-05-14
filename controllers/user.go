package controllers

import (
	"bytes"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"web/models"
	"web/modules/pager"
	"web/modules/passwords"
)

type UserController struct {
	BaseController
}
// Index 首页
func (c *UserController) Index() {
	c.Prepare()



	if c.User.Role != 0 {
		c.Abort("403")
	}

	c.Layout = ""
	c.TplName = "user/list.html"
	c.Data["MemberSelected"] = true

	pageIndex, _ := c.GetInt("page", 1)

	var users []models.User

	pageOptions := pager.PageOptions{
		TableName:  models.NewUser().TableName(),
		EnableFirstLastLink : true,
		CurrentPage : pageIndex,
		PageSize : 15,
		ParamName : "page",
		Conditions : " order by user_id desc",
	}



	//返回分页信息,
	//第一个:为返回的当前页面数据集合,ResultSet类型
	//第二个:生成的分页链接
	//第三个:返回总记录数
	//第四个:返回总页数
	totalItem, totalCount, rs, pageHtml := pager.GetPagerLinks(&pageOptions, c.Ctx)

	_,err := rs.QueryRows(&users)      //把当前页面的数据序列化进一个切片内

	if err != nil {
		logs.Error("",err.Error())
	}

	c.Data["lists"] = users
	c.Data["html"] = pageHtml
	c.Data["totalItem"] = totalItem
	c.Data["totalCount"] = totalCount

}


// My 个人中心
func (c *UserController) My(){
	c.Prepare()
	c.Layout = ""
	c.TplName = "user/edit.html"
	c.Data["MemberSelected"] = true

	user := c.User

	if c.Ctx.Input.IsPost() {
		password := c.GetString("password")
		status, _ := c.GetInt("status", 0)

		if password != "" {
			pass, _ := passwords.PasswordHash(password)
			user.Password = pass
		}
		if user.Role != 0 {
			user.Status = status
		}

		user.Email = c.GetString("email")
		user.Phone = c.GetString("phone")
		user.Avatar = c.GetString("avatar")

		if user.Avatar == "" {
			user.Avatar = "/static/images/headimgurl.jpg"
		}

		var result error

		if user.UserId > 0 {

			result = user.Update()
		} else {
			result = user.Add()
		}

		if result != nil {
			c.JsonResult(500, result.Error())
		}

		view, err := c.ExecuteViewPathTemplate("user/list_item.html", *user)
		if err != nil {
			logs.Error("", err.Error())
		}

		data := map[string]interface{}{
			"view" : view,
		}
		c.SetUser(*user)
		c.JsonResult(0, "ok", data)

	}

	c.Data["Model"] = user
	c.Data["IsSelf"] = true

}

// Edit 编辑信息
func (c *UserController) Edit() {
	c.TplName = "user/edit.html"

	c.Prepare()


	user_id ,_ := c.GetInt(":id")

	user := models.NewUser()

	if user_id > 0 {
		user.UserId = user_id
		if err := user.Find(); err != nil {
			c.ServerError("Data query error:" + err.Error())
		}
	}

	if c.Ctx.Input.IsPost() {
		password := c.GetString("password")
		account := c.GetString("account")
		status ,_ := c.GetInt("status",0)

		if user.UserId > 0 {
			if password != "" {
				pass, _ := passwords.PasswordHash(password)
				user.Password = pass
			}
			if user.Role != 0 {
				user.Status = status
			}
		} else {
			if account == ""{
				c.JsonResult(500,"Account is require.")
			}
			if password == "" {
				c.JsonResult(500,"Password is require.")
			}
			user.Role = 1
			user.Account = account
			user.Password = password
		}

		user.Email = c.GetString("email")
		user.Phone = c.GetString("phone")
		user.Avatar = c.GetString("avatar")

		if user.Avatar == "" {
			user.Avatar = "/static/images/headimgurl.jpg"
		}

		var result error

		if user.UserId > 0 {

			result = user.Update()
		} else {
			result = user.Add()
		}

		if result != nil {
			c.JsonResult(500, result.Error())
		}

		view, err := c.ExecuteViewPathTemplate("user/list_item.html", *user)
		if err != nil {
			logs.Error("", err.Error())
		}

		data := map[string]interface{}{
			"view" : view,
		}
		if user.UserId == c.User.UserId {
			c.SetUser(*user)
		}
		c.JsonResult(0, "ok", data)

	}


	c.Data["Model"] = user
	c.Data["IsSelf"] = false

	if user.UserId == c.User.UserId {
		c.Data["IsSelf"] = true
	}
}


// Delete 删除会员
func (c *UserController) Delete() {
	c.Prepare()

	if c.User.Role != 0 {
		c.Abort("403")
	}

	user_id ,err := c.GetInt(":id")

	if err != nil {
		logs.Error("",err.Error())
		c.JsonResult(500,"Parameter error.")
	}

	user := models.NewUser()
	user.UserId = user_id

	if err := user.Find();err != nil {
		logs.Error("",err.Error())
		c.JsonResult(500,"Data query error.")
	}

	if user.Role == 0 {
		c.JsonResult(500,"不能删除管理员用户")
	}
	if err := user.Delete();err != nil {
		logs.Error("",err.Error())
		c.JsonResult(500,"删除失败")
	}
	c.JsonResult(0,"ok")
}


// Upload 上传图片
func (c *UserController) Upload() {
	file,moreFile,err := c.GetFile("image-file")
	defer file.Close()

	if err != nil {
		logs.Error("",err.Error())
		c.JsonResult(500,"读取文件异常")
	}

	ext := filepath.Ext(moreFile.Filename)

	if !strings.EqualFold(ext,".png") && !strings.EqualFold(ext,".jpg") && !strings.EqualFold(ext,".gif") && !strings.EqualFold(ext,".jpeg")  {
		c.JsonResult(500,"不支持的图片格式")
	}


	x1 ,_ := strconv.ParseFloat(c.GetString("x"),10)
	y1 ,_ := strconv.ParseFloat(c.GetString("y"),10)
	w1 ,_ := strconv.ParseFloat(c.GetString("width"),10)
	h1 ,_ := strconv.ParseFloat(c.GetString("height"),10)

	x := int(x1)
	y := int(y1)
	width := int(w1)
	height := int(h1)

	fmt.Println(x,x1,y,y1)

	fileName := "avatar_" +  strconv.FormatInt(int64(time.Now().Nanosecond()), 16)

	filePath := "static/uploads/" + time.Now().Format("200601") + "/" + fileName + ext

	path := filepath.Dir(filePath)

	os.MkdirAll(path, os.ModePerm)

	err = c.SaveToFile("image-file",filePath)

	if err != nil {
		logs.Error("",err)
		c.JsonResult(500,"图片保存失败")
	}

	fileBytes,err := ioutil.ReadFile(filePath)

	if err != nil {
		logs.Error("",err)
		c.JsonResult(500,"图片保存失败")
	}

	buf := bytes.NewBuffer(fileBytes)

	m,_,err := image.Decode(buf)

	if err != nil{
		logs.Error("image.Decode => ",err)
		c.JsonResult(500,"图片解码失败")
	}


	var subImg image.Image

	if rgbImg,ok := m.(*image.YCbCr); ok {
		subImg = rgbImg.SubImage(image.Rect(x, y, x+width, y+height)).(*image.YCbCr) //图片裁剪x0 y0 x1 y1
	}else if rgbImg,ok := m.(*image.RGBA); ok {
		subImg = rgbImg.SubImage(image.Rect(x, y, x+width, y+height)).(*image.YCbCr) //图片裁剪x0 y0 x1 y1
	}else if rgbImg,ok := m.(*image.NRGBA); ok {
		subImg = rgbImg.SubImage(image.Rect(x, y, x+width, y+height)).(*image.YCbCr) //图片裁剪x0 y0 x1 y1
	} else {
		fmt.Println(m)
		c.JsonResult(500,"图片解码失败")
	}

	f, err := os.OpenFile("./" + filePath, os.O_SYNC | os.O_RDWR, 0666)

	if err != nil{
		c.JsonResult(500,"保存图片失败")
	}
	defer f.Close()

	if strings.EqualFold(ext,".jpg") || strings.EqualFold(ext,".jpeg"){

		err = jpeg.Encode(f,subImg,&jpeg.Options{ Quality : 100 })
	}else if strings.EqualFold(ext,".png") {
		err = png.Encode(f,subImg)
	}else if strings.EqualFold(ext,".gif") {
		err = gif.Encode(f,subImg,&gif.Options{ NumColors : 256})
	}
	if err != nil {
		logs.Error("图片剪切失败 => ",err.Error())
		c.JsonResult(500,"图片剪切失败")
	}


	if err != nil {
		logs.Error("保存文件失败 => ",err.Error())
		c.JsonResult(500,"保存文件失败")
	}
	url := "/" + filePath

	c.JsonResult(0,"ok",url)
}