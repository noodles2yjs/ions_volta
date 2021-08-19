package cars

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

type CarController struct {
	beego.Controller
}

func (u *CarController) IsActive() {

	is_active, _ := u.GetInt("status")
	id, _ := u.GetInt("id")
	o := orm.NewOrm()
	qs := o.QueryTable(new(auth.Car)).Filter("id", id)
	message_map := map[string]interface{}{}
	if is_active == 1 {
		qs.Update(orm.Params{
			"IsActive": 0,
		})
		ret := fmt.Sprintf("用户id:%d,停用成功", id)
		logs.Info(ret)
		message_map["msg"] = "停用成功"
	} else if is_active == 0 {
		qs.Update(orm.Params{
			"IsActive": 1,
		})
		ret := fmt.Sprintf("用户id:%d,启用成功", id)
		logs.Info(ret)
		message_map["msg"] = "启用成功"
	}

	u.Data["json"] = message_map
	u.ServeJSON()
}
func (c *CarController) List() {
	// 需要分页

	current_page, _ := c.GetInt("page", 1)
	// 查询
	kw := c.GetString("kw")
	pagesize := 2

	// 这是分页封装的方法，我会放在最下面(这个是重点) agricultures--对应要查询的models
	count, agricultures := CarListPageModel(current_page, pagesize, kw)

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
	c.Data["cars"] = agricultures
	c.Data["kw"] = kw
	c.TplName = "cars/car_list.html"

}

// 封装的列表分页方法
func CarListPageModel(page, pageSize int, kw string) (count int64, agricultures []auth.Car) {
	o := orm.NewOrm()

	qs := o.QueryTable(new(auth.Car))

	if kw != "" { // 有查询条件的
		// 总数
		count, _ = qs.Filter("IsDelete", 0).Filter("Name__contains", kw).Count()
		qs.Filter("IsDelete", 0).Filter("Name__contains", kw).Limit(pageSize).Offset((page - 1) * pageSize).RelatedSel().All(&agricultures)
	} else {
		count, _ = qs.Filter("IsDelete", 0).Count()
		qs.Filter("IsDelete", 0).Limit(pageSize).Offset((page - 1) * pageSize).RelatedSel().All(&agricultures)

	}

	return count, agricultures
}

func (c *CarController) ToAdd() {
	car_brand := []auth.CarBrand{}
	orm.NewOrm().QueryTable(new(auth.CarBrand)).Filter("IsDelete", 0).Filter("IsActive", 1).All(&car_brand)
	c.Data["car_brand"] = car_brand

	c.TplName = "cars/CarAdd.html"

}

func (c *CarController) DoAdd() {

	name := c.GetString("name")
	car_brand_id, _ := c.GetInt("carBrand")
	desc := c.GetString("desc")
	is_active, _ := c.GetInt("is_active")
	//时间的时区转换
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	timeStr2, _ := utils.ParseWithLocation("Asia/Shanghai", timeStr)

	car_brand := auth.CarBrand{Id: car_brand_id}
	o := orm.NewOrm()
	CarData := auth.Car{Name: name, Desc: desc, CarBand: &car_brand, IsActive: is_active, CreateTime: timeStr2}
	_, err := o.Insert(&CarData)
	message_map := map[string]interface{}{}
	if err != nil {
		ret1 := fmt.Sprintf("插入的数据信息:name:%s|desc:%s|is_active:%s", name, desc, is_active)
		ret := fmt.Sprintf("添加数据出错,错误信息为:%v", err)
		logs.Error(ret1)
		logs.Error(ret)
		message_map["code"] = 1001
		message_map["msg"] = "添加数据出错,请重新添加!"
	} else {
		ret := fmt.Sprintf("插入的数据信息:name:%s|desc:%s|is_active:%s", name, desc, is_active)
		logs.Info(ret)
		message_map["code"] = 200
		message_map["msg"] = "添加数据成功"

	}
	c.Data["json"] = message_map
	c.ServeJSON()
}

func (c *CarController) Delete() {
	id, _ := c.GetInt("id")
	o := orm.NewOrm()
	// 逻辑删除
	_, err := o.QueryTable(new(auth.Car)).Filter("id", id).Update(orm.Params{
		"IsDelete": 1,
	})

	message_map := map[string]interface{}{}
	if err != nil { //说明逻辑删除数据有问题
		ret1 := fmt.Sprintf("逻辑删除数据信息：车辆Id:%d", id)
		ret := fmt.Sprintf("逻辑删除数据出错,错误信息:%v", err)
		logs.Error(ret1)
		logs.Error(ret)
		message_map["code"] = 10001
		message_map["msg"] = "删除数据出错，请重新逻辑删除"

	} else {
		ret1 := fmt.Sprintf("删除数据信息成功：车辆Id:%d", id)
		logs.Info(ret1)
		message_map["code"] = 200
		message_map["msg"] = "删除成功"

	}
	c.Data["json"] = message_map
	c.ServeJSON()
}

func (c *CarController) ToUpdate() {
	id, _ := c.GetInt("id") // 车辆的id
	o := orm.NewOrm()
	qs := o.QueryTable(new(auth.CarBrand))
	// 通过id 查询 车辆所属品牌
	carBrand_sel := auth.CarBrand{}
	qs.Filter("Car__Id", id).One(&carBrand_sel)
	carBrands := []auth.CarBrand{}
	qs.Filter("IsDelete", 0).Filter("IsActive", 1).Exclude("id", carBrand_sel.Id).All(&carBrands)

	CarData := auth.Car{}
	o.QueryTable(new(auth.Car)).Filter("id", id).One(&CarData)
	c.Data["CarData"] = CarData
	c.Data["carBrand_sel"] = carBrand_sel
	c.Data["carBrands"] = carBrands
	c.TplName = "cars/Car_edit.html"

}
func (c *CarController) DoUpdate() {
	id, _ := c.GetInt("id")
	carBrandId, _ := c.GetInt("carBrandId")
	name := c.GetString("name")
	desc := c.GetString("desc")
	is_active, _ := c.GetInt("is_active")

	//carBrand_data:=auth.CarBrand{Id:carBrandId}
	o := orm.NewOrm()
	_, err := o.QueryTable(new(auth.Car)).Filter("id", id).Update(orm.Params{
		"Name":     name,
		"CarBand":  carBrandId,
		"Desc":     desc,
		"IsActive": is_active,
	})

	message_map := map[string]interface{}{}
	if err != nil {
		ret1 := fmt.Sprintf("数据信息为:name:%s|desc:%s|is_active:%d|", name, desc, is_active)
		ret := fmt.Sprintf("更新数据出错:%v", err)
		logs.Error(ret1)
		logs.Error(ret)
		message_map["code"] = 1001
		message_map["msg"] = "车辆更新出错,请重新更新!"
	} else {
		ret := fmt.Sprintf("更新数据成功,数据信息为:name:%s|desc:%s|is_active:%d|", name, desc, is_active)
		logs.Info(ret)
		message_map["code"] = 200
		message_map["msg"] = "车辆更新成功!"
	}
	c.Data["json"] = message_map
	c.ServeJSON()

}
func (c *CarController) MuliDelete() {

	ids := c.GetString("ids")
	id_arr := strings.Split(ids, ",")
	o := orm.NewOrm()
	qs := o.QueryTable(new(auth.Car))
	for _, v := range id_arr {
		id_int, _ := strconv.Atoi(v)
		fmt.Println(id_int)
		qs.Filter("id", id_int).Update(orm.Params{
			"IsDelete": 1,
		})
	}
	message_map := map[string]interface{}{}
	message_map["code"] = 200
	message_map["msg"] = "批量删除成功"

	ret := fmt.Sprintf("批量删除成功,车辆ids:%s", ids)
	logs.Info(ret)
	c.Data["json"] = message_map
	c.ServeJSON()

}
