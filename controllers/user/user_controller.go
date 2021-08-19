package user

import (
	"IOS_Volta/models/auth"
	"IOS_Volta/utils"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
	"time"
)

type UserController struct {
	beego.Controller
}

func (u *UserController) IsActive() {

	is_active_str := u.GetString("status")
	is_active := utils.StrToInt(is_active_str)
	//fmt.Println("============")
	//fmt.Println(is_active)
	id, _ := u.GetInt("id")
	//fmt.Println("============")
	//fmt.Println(id)
	o := orm.NewOrm()
	qs := o.QueryTable("sys_user").Filter("id", id)
	message_map := map[string]interface{}{}
	if is_active == 1 {
		qs.Update(orm.Params{
			"is_active": 0,
		})
		ret := fmt.Sprintf("用户id:%d,停用成功", id)
		logs.Info(ret)
		message_map["msg"] = "停用成功"
	} else if is_active == 0 {
		qs.Update(orm.Params{
			"is_active": 1,
		})
		ret := fmt.Sprintf("用户id:%d,启用成功", id)
		logs.Info(ret)
		message_map["msg"] = "启用成功"
	}

	u.Data["json"] = message_map
	u.ServeJSON()
}
func (u *UserController) List() {
	current_page, _ := u.GetInt("page", 1)
	// 查询
	kw := u.GetString("kw")
	pagesize := 7

	// 这是分页封装的方法，我会放在最下面(这个是重点) agricultures--对应要查询的models
	count, agricultures := AgricultureListModel(current_page, pagesize, kw)

	// 这是封装的方法，我会放在最下面
	num_pages := utils.GetPageNum(count, pagesize)

	// 这是封装的方法，我会放在最下面
	left_pages, right_pages, left_has_more, right_has_more := utils.Get_pagination_data(num_pages, current_page, 2)
	has_previous, has_next, previous_page_number, next_page_number := utils.HasNext(current_page, num_pages)

	u.Data["left_pages"] = left_pages         //左边页码数[]int
	u.Data["right_pages"] = right_pages       //右边页码数[]int
	u.Data["left_has_more"] = left_has_more   //左边是否有页码 bool
	u.Data["right_has_more"] = right_has_more //右边是否有页码 bool

	u.Data["has_previous"] = has_previous                 //是否有上一页 bool
	u.Data["has_next"] = has_next                         //是否有下一页 bool
	u.Data["previous_page_number"] = previous_page_number //上一页页码 int
	u.Data["next_page_number"] = next_page_number         //下一页页码 int

	u.Data["num_pages"] = num_pages //总页数
	u.Data["page"] = current_page   //当前页

	u.Data["count"] = count //总数量
	u.Data["users"] = agricultures
	u.Data["kw"] = kw
	u.TplName = "user/user_list.html"

}

// 封装的列表分页
func AgricultureListModel(page, pageSize int, kw string) (count int64, agricultures []auth.User) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(auth.User))
	if kw != "" { // 有查询条件的
		// 总数
		count, _ = qs.Filter("is_delete", 0).Filter("UserName__contains", kw).Count()
		qs.Filter("is_delete", 0).Filter("UserName__contains", kw).Limit(pageSize).Offset((page - 1) * pageSize).RelatedSel().All(&agricultures)
	} else {
		count, _ = qs.Filter("is_delete", 0).Count()
		qs.Filter("is_delete", 0).Limit(pageSize).Offset((page - 1) * pageSize).RelatedSel().All(&agricultures)

	}

	return count, agricultures
}

func (u *UserController) ToAdd() {
	u.TplName = "user/user_add.html"
}
func (u *UserController) DoAdd() {

	username := u.GetString("username")
	password := u.GetString("password")
	age, _ := u.GetInt("age")
	gender, _ := u.GetInt("gender")
	phone := u.GetString("phone")
	addr := u.GetString("addr")
	is_active, _ := u.GetInt("is_active")
	//时间的时区转换
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	timeStr2, _ := utils.ParseWithLocation("Asia/Shanghai", timeStr)
	new_password := utils.GetMd5Str(password)
	phone_int64, _ := strconv.ParseInt(phone, 10, 64)
	o := orm.NewOrm()
	user_data := auth.User{UserName: username, Password: new_password, Age: age, Gender: gender, Phone: phone_int64, Addr: addr, IsActive: is_active, CreateTime: timeStr2}
	_, err := o.Insert(&user_data)

	message_map := map[string]interface{}{}
	if err != nil { //说明插入数据有问题
		ret1 := fmt.Sprintf("插入数据信息：username:%s|md5_password:%s|age:%d|gender:%d|phone:%s|"+
			"addr:%s;is_active:%d", username, new_password, age, gender, phone, addr, is_active)
		ret := fmt.Sprintf("添加数据出错,错误信息:%v", err)
		logs.Error(ret1)
		logs.Error(ret)
		message_map["code"] = 10001
		message_map["msg"] = "添加数据出错，请重新添加"
		u.Data["json"] = message_map
	} else {
		ret1 := fmt.Sprintf("插入数据成功，数据信息：username:%s|md5_password:%s|age:%d|gender:%d|phone:%s|"+
			"addr:%s;is_active:%d", username, new_password, age, gender, phone, addr, is_active)
		logs.Info(ret1)
		message_map["code"] = 200
		message_map["msg"] = "添加成功"
		u.Data["json"] = message_map
	}
	u.ServeJSON()

}
func (u *UserController) ToUpdate() {
	id, _ := u.GetInt("id")
	userInfo := auth.User{}
	o := orm.NewOrm()
	err := o.QueryTable(new(auth.User)).Filter("id", id).One(&userInfo)
	var ret = ""
	if err != nil {
		ret = fmt.Sprintf("用户信息修改的ToUpdate()失败:%v", err)
		logs.Info(ret)
	} else {
		ret = fmt.Sprintf("用户信息修改的ToUpdate()成功,用户id:%d", id)
		logs.Info(ret)
		u.Data["user"] = userInfo
		u.TplName = "user/user_edit.html"

	}

}
func (u *UserController) DoUpdate() {
	user_id, _ := u.GetInt("user_id")
	username := u.GetString("username")
	password := u.GetString("password")
	age, _ := u.GetInt("age")
	gender, _ := u.GetInt("gender")
	phone := u.GetString("phone")
	addr := u.GetString("addr")
	is_active, _ := u.GetInt("is_active")

	//时间的时区转换
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	timeStr2, _ := utils.ParseWithLocation("Asia/Shanghai", timeStr)
	new_password := utils.GetMd5Str(password)
	phone_int64, _ := strconv.ParseInt(phone, 10, 64)
	o := orm.NewOrm()
	_, err := o.QueryTable(new(auth.User)).Filter("id", user_id).Update(orm.Params{
		"UserName":   username,
		"Password":   new_password,
		"Age":        age,
		"Gender":     gender,
		"Phone":      phone_int64,
		"Addr":       addr,
		"IsActive":   is_active,
		"CreateTime": timeStr2,
	})

	message_map := map[string]interface{}{}
	if err != nil { //说明更新数据有问题
		ret1 := fmt.Sprintf("更新数据信息：username:%s|md5_password:%s|age:%d|gender:%d|phone:%s|"+
			"addr:%s;is_active:%d", username, new_password, age, gender, phone, addr, is_active)
		ret := fmt.Sprintf("更新数据出错,错误信息:%v", err)
		logs.Error(ret1)
		logs.Error(ret)
		message_map["code"] = 10001
		message_map["msg"] = "更新数据出错，请重新更新"
		u.Data["json"] = message_map
	} else {
		ret1 := fmt.Sprintf("更新数据成功，数据信息：username:%s|md5_password:%s|age:%d|gender:%d|phone:%s|"+
			"addr:%s;is_active:%d", username, new_password, age, gender, phone, addr, is_active)
		logs.Info(ret1)
		message_map["code"] = 200
		message_map["msg"] = "更新成功"
		u.Data["json"] = message_map
	}

	u.ServeJSON()

}
func (u *UserController) Delete() {
	id, _ := u.GetInt("id")

	o := orm.NewOrm()
	//逻辑性删除
	_, err := o.QueryTable("sys_user").Filter("id", id).Update(orm.Params{
		"is_delete": 1,
	})
	message_map := map[string]interface{}{}
	if err != nil { //说明逻辑删除数据有问题
		ret1 := fmt.Sprintf("逻辑删除数据信息：用户Id:%d", id)
		ret := fmt.Sprintf("逻辑删除数据出错,错误信息:%v", err)
		logs.Error(ret1)
		logs.Error(ret)
		message_map["code"] = 10001
		message_map["msg"] = "逻辑删除数据出错，请重新逻辑删除"
		u.Data["json"] = message_map
	} else {
		ret1 := fmt.Sprintf("逻辑删除数据信息成功：用户Id:%d", id)
		logs.Info(ret1)
		message_map["code"] = 200
		message_map["msg"] = "逻辑删除成功"
		u.Data["json"] = message_map
	}

	u.ServeJSON()

}
func (u *UserController) ResetPassword() {
	id, _ := u.GetInt("id")
	o := orm.NewOrm()
	//重置密码
	new_pwd := utils.GetMd5Str("123456")
	_, err := o.QueryTable("sys_user").Filter("id", id).Update(orm.Params{
		"password": new_pwd,
	})
	message_map := map[string]interface{}{}
	if err != nil { //说明重置密码数据有问题
		ret1 := fmt.Sprintf("重置密码数据信息：用户Id:%d", id)
		ret := fmt.Sprintf("重置密码数据出错,错误信息:%v", err)
		logs.Error(ret1)
		logs.Error(ret)
		message_map["code"] = 10001
		message_map["msg"] = "重置密码数据出错，请重新重置密码"
		u.Data["json"] = message_map
	} else {
		ret1 := fmt.Sprintf("重置密码数据信息成功：用户Id:%d", id)
		logs.Info(ret1)
		message_map["code"] = 200
		message_map["msg"] = "重置密码成功,重置后的默认密码为123456!"
		u.Data["json"] = message_map
	}

	u.ServeJSON()

}

func (u *UserController) MultiDelete() {

	ids := u.GetString("ids")
	id_arr := strings.Split(ids, ",")
	o := orm.NewOrm()
	qs := o.QueryTable("sys_user")
	for _, v := range id_arr {
		id_int := utils.StrToInt(v)
		qs.Filter("id", id_int).Update(orm.Params{
			"is_delete": 1,
		})

	}

	ret := fmt.Sprintf("批量删除成功，用户ids:%d", ids)
	logs.Info(ret)

	u.Data["json"] = map[string]interface{}{"code": 200, "msg": "批量删除成功"}
	u.ServeJSON()

}
