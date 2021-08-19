package auth

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type Auth struct {
	Id         int       `orm:"pk;auto"`
	AuthName   string    `orm:"description(权限名称);size(64)"`
	UrlFor     string    `orm:"description(url反转);size(255)"`
	Pid        int       `orm:"description(父节点id)"`
	Desc       string    `orm:"description(描述);size(255)"`
	CreateTime time.Time `orm:"type(datetime);auto_now;description(创建时间)"`
	IsActive   int       `orm:"column(is_active);description(1启用，0停用)"`
	IsDelete   int       `orm:"columns(is_delete);description(1删除，0未删除)"`
	Weight     int       `orm:"description(权重，数值越大，权重越大)"`
	Role       []*Role   `orm:"reverse(many)"`
}

type Role struct {
	Id         int       `orm:"pk;auto"`
	RoleName   string    `orm:"size(64)"`
	Desc       string    `orm:"size(255)"`
	IsActive   int       `orm:"column(is_active)"`
	IsDelete   int       `orm:"column(is_delete)"`
	CreateTime time.Time `orm:"type(datetime);auto_now;description(创建时间)"`
	Auth       []*Auth   `orm:"rel(m2m)"`
	User       []*User   `orm:"rel(m2m)"`
}

type User struct {
	Id             int               `orm:"pk;auto"`
	CardId         string            `orm:"size(64);column(card_id);description(员工工号)"`
	UserName       string            `orm:"column(user_name);size(64);description(用户名)"`
	Password       string            `orm:"size(32);description(密码)"`
	Age            int               `orm:"null;description(年龄)"`
	Gender         int               `orm:"null;description(1:男,2:女,3:未知)"`
	Phone          int64             `orm:"null;description(电话号码)"`
	Addr           string            `orm:"null;size(255);description(地址)"`
	Role           []*Role           `orm:"reverse(many)"`
	MessageNotify  []*MessageNotify  `orm:"reverse(many)"`
	CarApplication []*CarApplication `orm:"reverse(many)"`
	IsActive       int               `orm:"description(1启用，0停用);default(1)"`
	IsDelete       int               `orm:"description(1删除，0未删除);default(0)"`
	CreateTime     time.Time         `orm:"auto_now;type(datetime);description(创建时间);null"`
}

//车辆品牌
type CarBrand struct {
	Id         int       `orm:"pk;auto;"`
	Name       string    `orm:"column(Name);size(64);desciption(品牌名字);column(Name)"`
	Car        []*Car    `orm:"reverse(many)"`
	Desc       string    `orm:"column(Desc);size(255);desciption(品牌描述);column(Desc)"`
	IsActive   int       `orm:"column(IsActive);default(1);desciption(1:启用,0:停用)"`
	IsDelete   int       `orm:"column(IsDelete);default(0);desciption(1:删除,0:未删除)"`
	CreateTime time.Time `orm:"column(CreateTime);type(datetime);auto_now;description(创建时间)"`
}

//车辆列表
type Car struct {
	Id             int               `orm:"pk;auto;"`
	Name           string            `orm:"column(Name);size(64);desciption(车辆名称);column(Name)"`
	Desc           string            `orm:"column(Desc);size(255);desciption(品牌描述);column(Desc)"`
	CarBand        *CarBrand         `orm:"rel(fk);"`
	CarApplication []*CarApplication `orm:"reverse(many)"`
	Status         int               `orm:"column(Status);description(1:可借,2:不可借);default(1)"`
	IsActive       int               `orm:"column(IsActive);default(1);desciption(1:启用,0:停用)"`
	IsDelete       int               `orm:"column(IsDelete);default(0);desciption(1:删除,0:未删除)"`
	CreateTime     time.Time         `orm:"column(CreateTime);type(datetime);auto_now;description(创建时间)"`
}

// 车辆申请
type CarApplication struct {
	Id           int       `orm:"pk;auto;"`
	User         *User     `orm:"rel(fk);"`
	Car          *Car      `orm:"rel(fk);"`
	Reason       string    `orm:"description(申请理由);size(255);column(Reason)"`
	Destination  string    `orm:"description(目的地);size(64);column(Destination)"`
	ReturnDate   time.Time `orm:"type(date);auto_now;description(归还日期);column(ReturnDate)"`
	ReturnStatus int       `orm:"descrption(1:已归还，0：未归还);default(0);column(ReturnStatus)"`
	AuditStatus  int       `orm:"description(1:同意，2:未同意，3:未审批);default(3);column(AuditStatus)"`
	AuditOption  string    `orm:"description(审批意见);size(255);column(AuditOption)"`
	NotifyTag    int       `orm:"description(1:已发送通知，0：未发送通知);default(0);column(NotifyTag)"`
	IsActive     int       `orm:"column(IsActive);default(1);desciption(1:启用,0:停用)"`
	IsDelete     int       `orm:"column(IsDelete);default(0);desciption(1:删除,0:未删除)"`
	CreateTime   time.Time `orm:"column(CreateTime);type(datetime);auto_now;description(创建时间)"`
}

// 消息通知
type MessageNotify struct {
	Id      int    `orm:"pk;auto"`
	Flag    int    `orm:"column(Flag);description(1:车辆逾期，2:所有通知);default(1)"`
	Title   string `orm:"column(Title);size(64);description(消息标题)"`
	Content string `orm:"column(Content);type(text);description(消息内容)"`
	User    *User  `orm:"rel(fk);desription(用户外键)"`
	ReadTag int    `orm:"column(ReadTag);description(1:已读，0:未读)"`
}

func (u *User) TableName() string {
	return "sys_user"
}

func (r *Auth) TableName() string {
	return "sys_auth"

}

func (r *Role) TableName() string {
	return "sys_role"

}

func init() {
	orm.RegisterModel(
		new(Auth),
		new(Role),
		new(User),
		new(CarApplication),
		new(Car),
		new(CarBrand),
		new(MessageNotify),
	)

}
