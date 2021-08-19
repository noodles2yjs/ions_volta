package auth

import (
	"IOS_Volta/models/auth"
	"IOS_Volta/utils"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"time"
)

type AuthController struct {
	beego.Controller
}

func (a *AuthController) IsActive() {

}
func (u *AuthController) List() {

	current_page, _ := u.GetInt("page", 1)
	// 查询
	kw := u.GetString("kw")
	pagesize := 5

	// 这是分页封装的方法，我会放在最下面(这个是重点) agricultures--对应要查询的models
	count, agricultures := ListPageAuthModel(current_page, pagesize, kw)

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
	u.Data["auths"] = agricultures
	u.Data["kw"] = kw
	u.TplName = "auth/auth_list.html"

}

// 封装的列表分页方法
func ListPageAuthModel(page, pageSize int, kw string) (count int64, agricultures []auth.Auth) {
	o := orm.NewOrm()

	qs := o.QueryTable(new(auth.Auth))

	if kw != "" { // 有查询条件的
		// 总数
		count, _ = qs.Filter("is_delete", 0).Filter("AuthName__contains", kw).Count()
		qs.Filter("is_delete", 0).Filter("AuthName__contains", kw).Limit(pageSize).Offset((page - 1) * pageSize).RelatedSel().All(&agricultures)
	} else {
		count, _ = qs.Filter("is_delete", 0).Count()
		qs.Filter("is_delete", 0).Limit(pageSize).Offset((page - 1) * pageSize).RelatedSel().All(&agricultures)

	}

	return count, agricultures
}

func (a *AuthController) ToAdd() {
	o := orm.NewOrm()
	auths := []auth.Auth{}
	_, err := o.QueryTable("sys_auth").Filter("is_delete", 0).All(&auths)
	if err != nil {
		ret1 := fmt.Sprintf("查询数据出错,出错信息:%v", err)
		logs.Error(ret1)
	}

	a.Data["auths"] = auths
	a.TplName = "auth/auth_add.html"

}
func (a *AuthController) DoAdd() {
	auth_parent_id, _ := a.GetInt("auth_parent_id")
	auth_name := a.GetString("auth_name")
	auth_url := a.GetString("auth_url")
	auth_desc := a.GetString("auth_desc")
	is_active, _ := a.GetInt("is_active")
	auth_weight, _ := a.GetInt("auth_weight")
	//时间的时区转换
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	timeStr2, _ := utils.ParseWithLocation("Asia/Shanghai", timeStr)
	auth_data := auth.Auth{AuthName: auth_name, UrlFor: auth_url, Pid: auth_parent_id, Desc: auth_desc, IsActive: is_active, Weight: auth_weight, CreateTime: timeStr2}
	o := orm.NewOrm()
	_, err := o.Insert(&auth_data)
	if err != nil {
		ret := fmt.Sprint("添加权限出错:出错信息为:%v", err)
		logs.Error(ret)
	}

	a.Data["json"] = map[string]interface{}{"code": 200, ",msg": "添加权限成功!"}
	a.ServeJSON()

}
func (a *AuthController) ToUpdate() {

}
func (a *AuthController) DoUpdate() {

}
func (a *AuthController) Delete() {

}
func (a *AuthController) MultiDelete() {

}
