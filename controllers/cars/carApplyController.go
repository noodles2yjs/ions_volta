package cars

import (
	"IOS_Volta/models/auth"
	"IOS_Volta/utils"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type CarApplyController struct {
	beego.Controller
}

func (c *CarApplyController) List() {

	// 需要分页

	current_page, _ := c.GetInt("page", 1)
	// 查询
	kw := c.GetString("kw")
	pagesize := 2

	// 这是分页封装的方法，我会放在最下面(这个是重点) agricultures--对应要查询的models
	count, agricultures := CarApplyListPageModel(current_page, pagesize, kw)

	// 这是封装的方法，我会放在最下面
	num_pages := utils.GetPageNum(count, pagesize)

	// 这是封装的方法，我会放在最下面
	left_pages, right_pages, left_has_more, right_has_more := utils.Get_pagination_data(num_pages, current_page, 2)
	has_previous, has_next, previous_page_number, next_page_number := utils.HasNext(current_page, num_pages)

	c.Data["left_pages"] = left_pages         //左边页码数[]int
	c.Data["right_pages"] = right_pages       //右边页码数[]int
	c.Data["left_has_more"] = left_has_more   //左边是否有页码 bool
	c.Data["right_has_more"] = right_has_more //右边是否有页码 bool

	c.Data["has_previous"] = has_previous                 //是否有上一页 bool
	c.Data["has_next"] = has_next                         //是否有下一页 bool
	c.Data["previous_page_number"] = previous_page_number //上一页页码 int
	c.Data["next_page_number"] = next_page_number         //下一页页码 int

	c.Data["num_pages"] = num_pages //总页数
	c.Data["page"] = current_page   //当前页

	c.Data["count"] = count //总数量
	c.Data["cars_data"] = agricultures
	c.Data["kw"] = kw
	c.TplName = "cars/carApply_list.html"

}

func CarApplyListPageModel(page, pageSize int, kw string) (count int64, agricultures []auth.Car) {
	o := orm.NewOrm()

	qs := o.QueryTable(new(auth.Car))

	if kw != "" { // 有查询条件的
		// 总数
		count, _ = qs.Filter("IsDelete", 0).Filter("Name__contains", kw).Count()
	} else {
		qs.Filter("IsDelete", 0).Filter("Name__contains", kw).Limit(pageSize).Offset((page - 1) * pageSize).RelatedSel().All(&agricultures)
		count, _ = qs.Filter("IsDelete", 0).Count()
		qs.Filter("IsDelete", 0).Limit(pageSize).Offset((page - 1) * pageSize).RelatedSel().All(&agricultures)

	}

	return count, agricultures
}

func (c *CarApplyController) ToApply() {
	id, _ := c.GetInt("id")
	c.Data["id"] = id
	c.TplName = "cars/carApply.html"
}

func (c *CarApplyController) DoApply() {
	reason := c.GetString("reason")
	destination := c.GetString("destination")
	return_date, _ := time.Parse("2006-01-02", c.GetString("return_date"))
	cars_id, _ := c.GetInt("cars_id")
	uid := c.GetSession("id").(int)

	o := orm.NewOrm()
	cars := auth.Car{Id: cars_id}
	user_data := auth.User{Id: uid}
	cars_apply := auth.CarApplication{
		User:         &user_data,
		Car:          &cars,
		Reason:       reason,
		Destination:  destination,
		ReturnDate:   return_date,
		ReturnStatus: 0,
		AuditStatus:  3,
		IsActive:     1,
	}
	_, err := o.Insert(&cars_apply)
	o.QueryTable(new(auth.Car)).Filter("id", cars_id).Update(orm.Params{
		"Status": 2,
	})

	message_map := map[string]interface{}{}
	if err != nil {
		message_map["code"] = 10001
		message_map["msg"] = "申请失败"
	}

	message_map["code"] = 200
	message_map["msg"] = "申请成功"

	c.Data["json"] = message_map
	c.ServeJSON()

}
func (c *CarApplyController) MyApply() {

	// 需要分页

	current_page, _ := c.GetInt("page", 1)
	// 查询
	kw := c.GetString("kw")
	pagesize := 2
	uid := c.GetSession("id").(int)

	// 这是分页封装的方法，我会放在最下面(这个是重点) agricultures--对应要查询的models
	count, agricultures := MyApplyListPageModel(current_page, pagesize, kw, uid)

	// 这是封装的方法，我会放在最下面
	num_pages := utils.GetPageNum(count, pagesize)

	// 这是封装的方法，我会放在最下面
	left_pages, right_pages, left_has_more, right_has_more := utils.Get_pagination_data(num_pages, current_page, 2)
	has_previous, has_next, previous_page_number, next_page_number := utils.HasNext(current_page, num_pages)

	c.Data["left_pages"] = left_pages         //左边页码数[]int
	c.Data["right_pages"] = right_pages       //右边页码数[]int
	c.Data["left_has_more"] = left_has_more   //左边是否有页码 bool
	c.Data["right_has_more"] = right_has_more //右边是否有页码 bool

	c.Data["has_previous"] = has_previous                 //是否有上一页 bool
	c.Data["has_next"] = has_next                         //是否有下一页 bool
	c.Data["previous_page_number"] = previous_page_number //上一页页码 int
	c.Data["next_page_number"] = next_page_number         //下一页页码 int

	c.Data["num_pages"] = num_pages //总页数
	c.Data["page"] = current_page   //当前页

	c.Data["count"] = count //总数量
	c.Data["cars_apply"] = agricultures
	c.Data["kw"] = kw
	c.TplName = "cars/myApply_list.html"

}

func MyApplyListPageModel(page, pageSize int, kw string, uid int) (count int64, agricultures []auth.CarApplication) {
	o := orm.NewOrm()

	qs := o.QueryTable(new(auth.CarApplication))

	if kw != "" { // 有查询条件的
		// 总数
		qs.Filter("IsDelete", 0).Filter("Car__Name__contains", kw).Filter("user_id", uid).Limit(pageSize).Offset((page - 1) * pageSize).RelatedSel().All(&agricultures)
		count, _ = qs.Filter("IsDelete", 0).Filter("Car__Name__contains", kw).Filter("user_id", uid).Count()
	} else {
		qs.Filter("IsDelete", 0).Filter("user_id", uid).Limit(pageSize).Offset((page - 1) * pageSize).RelatedSel().All(&agricultures)
		count, _ = qs.Filter("IsDelete", 0).Filter("user_id", uid).Count()

	}

	return count, agricultures
}

func (c *CarApplyController) AuditApply() {

	// 需要分页
	current_page, _ := c.GetInt("page", 1)
	// 查询
	kw := c.GetString("kw")
	pagesize := 2

	// 这是分页封装的方法，我会放在最下面(这个是重点) agricultures--对应要查询的models
	count, agricultures := AuditApplyListPageModel(current_page, pagesize, kw)

	// 这是封装的方法，我会放在最下面
	num_pages := utils.GetPageNum(count, pagesize)

	// 这是封装的方法，我会放在最下面
	left_pages, right_pages, left_has_more, right_has_more := utils.Get_pagination_data(num_pages, current_page, 2)
	has_previous, has_next, previous_page_number, next_page_number := utils.HasNext(current_page, num_pages)

	c.Data["left_pages"] = left_pages         //左边页码数[]int
	c.Data["right_pages"] = right_pages       //右边页码数[]int
	c.Data["left_has_more"] = left_has_more   //左边是否有页码 bool
	c.Data["right_has_more"] = right_has_more //右边是否有页码 bool

	c.Data["has_previous"] = has_previous                 //是否有上一页 bool
	c.Data["has_next"] = has_next                         //是否有下一页 bool
	c.Data["previous_page_number"] = previous_page_number //上一页页码 int
	c.Data["next_page_number"] = next_page_number         //下一页页码 int

	c.Data["num_pages"] = num_pages //总页数
	c.Data["page"] = current_page   //当前页

	c.Data["count"] = count //总数量
	c.Data["cars_apply"] = agricultures
	c.Data["kw"] = kw
	c.TplName = "cars/auditApply_list.html"

}

func AuditApplyListPageModel(page, pageSize int, kw string) (count int64, agricultures []auth.CarApplication) {

	o := orm.NewOrm()

	qs := o.QueryTable(new(auth.CarApplication))

	if kw != "" { // 有查询条件的
		// 总数
		count, _ = qs.Filter("IsDelete", 0).Filter("Car__Name__contains", kw).Count()
		qs.Filter("IsDelete", 0).Filter("Car__Name__contains", kw).Limit(pageSize).Offset((page - 1) * pageSize).RelatedSel().All(&agricultures)
	} else {
		count, _ = qs.Filter("IsDelete", 0).Count()
		qs.Filter("IsDelete", 0).Limit(pageSize).Offset((page - 1) * pageSize).RelatedSel().All(&agricultures)

	}

	return count, agricultures

}

func (c *CarApplyController) ToAuditApply() {
	id, _ := c.GetInt("id")
	o := orm.NewOrm()
	qs := o.QueryTable(new(auth.CarApplication))
	cars_apply := auth.CarApplication{}
	qs.Filter("id", id).One(&cars_apply)
	c.Data["cars_apply"] = cars_apply
	c.TplName = "cars/auditApply.html"

}

func (c *CarApplyController) DoAuditApply() {
	id, _ := c.GetInt("id")
	option := c.GetString("option")
	audit_status, _ := c.GetInt("audit_status")

	o := orm.NewOrm()
	qs := o.QueryTable(new(auth.CarApplication))
	_, err := qs.Filter("id", id).Update(orm.Params{
		"AuditOption": option,
		"AuditStatus": audit_status,
	})
	message_map := map[string]interface{}{}
	if err != nil {
		message_map["code"] = 10001
		message_map["msg"] = "审核失败"
	}
	message_map["code"] = 200
	message_map["msg"] = "审核成功"
	c.Data["json"] = message_map
	c.ServeJSON()
}

func (c *CarApplyController) DoReturn() {
	id, _ := c.GetInt("id")
	o := orm.NewOrm()
	qs := o.QueryTable(new(auth.CarApplication)).Filter("id", id)
	qs.Update(orm.Params{
		"ReturnStatus": 1,
	})
	car_apply := auth.CarApplication{}
	qs.One(&car_apply)
	o.QueryTable(new(auth.Car)).Filter("id", car_apply.Car.Id).Update(orm.Params{
		"Status": 1,
	})

	c.Redirect(beego.URLFor("CarApplyController.MyApply"), 302)
	// c.Redirect("CarApplyController.MyApply",302) -此方法无法重定向

}
