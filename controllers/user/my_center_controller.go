package user

import (
	"IOS_Volta/models/auth"
	"IOS_Volta/utils"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strings"
	"time"
)

type MyCenterController struct {
	beego.Controller
}

func (t *MyCenterController) Get() {

	user_id := t.GetSession("id")
	o := orm.NewOrm()
	userData := auth.User{}
	o.QueryTable(new(auth.User)).Filter("Id", user_id).One(&userData)
	t.Data["user"] = userData
	t.TplName = "user/my_center_edit.html"
}

func (t *MyCenterController) Post() {

	uid, _ := t.GetInt("uid")
	username := t.GetString("username")
	old_pwd := t.GetString("old_pwd")
	old_pwd = utils.GetMd5Str(old_pwd)
	new_pwd := t.GetString("new_pwd")
	new_pwd = utils.GetMd5Str(new_pwd)
	age, _ := t.GetInt("age")
	gender, _ := t.GetInt("gender")
	phone := utils.StrToInt(t.GetString("phone"))

	addr := t.GetString("addr")
	addr = strings.TrimSpace(addr)
	is_active, _ := t.GetInt("is_active")
	//时间的时区转换
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	timeStr2, _ := utils.ParseWithLocation("Asia/Shanghai", timeStr)

	o := orm.NewOrm()
	user := auth.User{}

	qs := o.QueryTable(new(auth.User)).Filter("id", uid)
	qs.One(&user)
	message_map := map[string]interface{}{}
	if old_pwd != user.Password {
		message_map["code"] = 1001
		message_map["msg"] = "原始密码错误"

	} else {
		qs.Update(orm.Params{
			"username":   username,
			"Password":   new_pwd,
			"age":        age,
			"gender":     gender,
			"phone":      phone,
			"addr":       addr,
			"is_active":  is_active,
			"CreateTime": timeStr2,
		})
		message_map["code"] = 200
		message_map["msg"] = "修改成功!"
	}

	t.Data["json"] = message_map
	t.ServeJSON()

}
