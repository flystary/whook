package models

import (
	"errors"
	"github.com/beego/beego/v2/client/orm"
	"time"
	"web/modules/hash"
)

// WebHook 对象
type WebHook struct {
	WebHookId 		int			`orm:"pk;auto;unique;column(web_hook_id)" json:"web_hook_id"`
	RepositoryName 	string		`orm:"size(255);column(repo_name)" json:"repository_name"`
	BranchName 		string		`orm:"size(255);column(branch_name)" json:"branch_name"`
	Tag 			string		`orm:"size(1000);column(tag)" json:"tag"`
	Shell 			string		`orm:"size(1000);column(shell)" json:"shell"`
	Status 			int			`orm:"type(int);column(status);default(0)" json:"status"`
	Key 			string		`orm:"size(255);column(key);unique" json:"key"`
	Secure 			string		`orm:"size(255);column(secure)" json:"secure"`
	LastExecTime 	time.Time	`orm:"type(datetime);column(last_exec_time);null" json:"last_exec_time"`
	CreateTime 		time.Time	`orm:"type(datetime);column(create_time);auto_now_add" json:"create_time"`
	HookType 		string 		`orm:"column(hook_type);size(50)" json:"hook_type"`
	CreateAt 		int			`orm:"type(int);column(create_at)"`
}

// TableName 获取对应数据库表名
func (m *WebHook) TableName() string {
	return "webhooks"
}

// TableEngine 获取数据使用的引擎
func (m *WebHook) TableEngine() string {
	return "INNODB"
}

// NewWebHook 新建 WebHook 对象
func NewWebHook() *WebHook {
	return &WebHook{}
}

// Find 查找
func (m *WebHook) Find() error {

	if m.WebHookId <= 0 {
		return ErrInvalidParameter
	}

	o := orm.NewOrm()


	if err := o.Read(m) ;err != nil {
		return err
	}
	return nil
}

// DeleteMulti 批量删除
func (m *WebHook) DeleteMulti (id ...int) error {
	if len(id) > 0 {
		o := orm.NewOrm()
		ids := make([]int,len(id))
		params := ""

		for i := 0;i<len(id);i++ {
			ids[i] = id[i]
			params = params + ",?"
		}
		_,err := o.Raw("DELETE webhooks WHERE web_hook_id IN ("+ params[1:] +")",ids).Exec()
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("invalid parameter")
}

// Delete 删除一条
func (m *WebHook) Delete() error {
	o := orm.NewOrm()
	_,err := o.Delete(m)

	return err
}

// FindByKey 根据Key查找
func (m *WebHook) FindByKey(key string) error {
	o := orm.NewOrm()

	if err := o.QueryTable(m.TableName()).Filter("key",key).One(m);err != nil {
		return err
	}
	return nil
}

// Save 添加或更新
func (m *WebHook) Save() error {
	o := orm.NewOrm()
	var err error

	if m.WebHookId > 0 {
		_,err = o.Update(m)
	}else{
		key := time.Now().String() + m.RepositoryName + m.BranchName

		m.Key = hash.Md5(key)

		m.Secure = hash.Md5(key + key + time.Now().String())

		_,err = o.Insert(m)
	}

	return err
}