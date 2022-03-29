package models

import (
	"github.com/beego/beego/v2/client/orm"
	"time"
	"web/modules/passwords"
)

// User 会员信息.
type User struct {
	UserId        int       `orm:"pk;auto;unique;column(user_id)"`
	Account       string    `orm:"size(100);unique;column(account)"`
	Password      string    `orm:"size(1000);column(password)"`
	Email         string    `orm:"size(255);column(email);null;default(null)"`
	Phone         string    `orm:"size(255);column(phone);null;default(null)"`
	Avatar        string    `orm:"size(1000);column(avatar)"`
	Role          int       `orm:"column(role);type(int);default(1)"`   //用户角色：0 管理员/1 普通用户
	Status        int       `orm:"column(status);type(int);default(0)"` //用户状态：0 正常/1 禁用
	CreateTime    time.Time `orm:"type(datetime);column(create_time);auto_now_add"`
	CreateAt      int       `orm:"type(int);column(create_at)"`
	LastLoginTime time.Time `orm:"type(datetime);column(last_login_time);null"`
}

// TableName 获取对应数据库表名.
func (u *User) TableName() string {
	return "users"
}

// TableEngine 获取数据使用的引擎.
func (u *User) TableEngine() string {
	return "INNODB"
}

// NewUser 获取新的用户信息对象.
func NewUser() *User {
	return &User{}
}

// Find 根据用户ID查找用户.
func (u *User) Find() error {
	o := orm.NewOrm()

	err := o.Read(u)

	if err == orm.ErrNoRows {
		return ErrMemberNoExist
	}

	return nil
}

// Login 用户登录.
func (u *User) Login(account string, password string) (*User, error) {
	o := orm.NewOrm()

	user := &User{}

	err := o.QueryTable(u.TableName()).Filter("account", account).Filter("status", 0).One(user)

	if err != nil {
		return user, ErrMemberNoExist
	}

	ok, err := passwords.PasswordVerify(user.Password, password)

	if ok && err == nil {
		return user, nil
	}

	return user, ErrorMemberPasswordError
}

// Add 添加一个用户.
func (u *User) Add() error {
	o := orm.NewOrm()

	hash, err := passwords.PasswordHash(u.Password)

	if err != nil {
		return err
	}

	u.Password = hash

	_, err = o.Insert(u)

	if err != nil {
		return err
	}
	return nil
}

// Update 更新用户信息.
func (u *User) Update(cols ...string) error {
	o := orm.NewOrm()

	if _, err := o.Update(u, cols...); err != nil {
		return err
	}
	return nil
}

// Delete 删除一个用户.
func (u *User) Delete() error {
	o := orm.NewOrm()

	if _, err := o.Delete(u); err != nil {
		return err
	}
	return nil
}
