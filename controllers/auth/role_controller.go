package auth

import (
	"IOS_Volta/models/auth"
	"IOS_Volta/utils"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strings"
)

type RoleController struct {
	beego.Controller
}

func (u *RoleController) List() {

	current_page, _ := u.GetInt("page", 1)
	// 查询
	kw := u.GetString("kw")
	pagesize := 5

	// 这是分页封装的方法，我会放在最下面(这个是重点) agricultures--对应要查询的models
	count, agricultures := ListPageRoleModel(current_page, pagesize, kw)

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
	u.Data["roles"] = agricultures
	u.Data["kw"] = kw
	u.TplName = "auth/role_list.html"
}

// 封装的列表分页方法
func ListPageRoleModel(page, pageSize int, kw string) (count int64, agricultures []auth.Role) {
	o := orm.NewOrm()

	qs := o.QueryTable(new(auth.Role))

	if kw != "" { // 有查询条件的
		// 总数
		count, _ = qs.Filter("is_delete", 0).Filter("AuthName__contains", kw).Count()
		qs.Filter("is_delete", 0).Filter("RoleName__contains", kw).Limit(pageSize).Offset((page - 1) * pageSize).RelatedSel().All(&agricultures)
	} else {
		count, _ = qs.Filter("is_delete", 0).Count()
		qs.Filter("is_delete", 0).Limit(pageSize).Offset((page - 1) * pageSize).RelatedSel().All(&agricultures)

	}

	return count, agricultures
}

func (r *RoleController) ToAdd() {

	r.TplName = "auth/role_add.html"
}

func (r *RoleController) DoAdd() {

}

func (r *RoleController) ToRoleAuth() {

	role_id, _ := r.GetInt("role_id")
	o := orm.NewOrm()
	qs := o.QueryTable("sys_role")
	role := auth.Role{}
	qs.Filter("id", role_id).One(&role)
	r.Data["role"] = role
	r.TplName = "auth/role-auth-add.html"

}

func (r *RoleController) DoRoleAuth() {
	role_id, _ := r.GetInt("role_id")
	auth_ids := r.GetString("auth_ids")
	new_auth_ids := auth_ids[0 : len(auth_ids)-1]
	auth_id_arry := strings.Split(new_auth_ids, ",")

	role := auth.Role{Id: role_id}
	o := orm.NewOrm()
	//  // 先清空之前数据库里的添加的数据
	m2m := o.QueryM2M(&role, "Auth")
	m2m.Clear()

	//在添加现在选中放在数据库中
	for _, auth_id := range auth_id_arry {
		auth_id_int := utils.StrToInt(auth_id)
		auth_data := auth.Auth{Id: auth_id_int}
		m2m.Add(&auth_data)
	}
	message_map := map[string]interface{}{}
	message_map["code"] = 200
	message_map["msg"] = "添加成功"
	r.Data["json"] = message_map
	r.ServeJSON()
}

func (r *RoleController) ToRoleUser() {
	role_id, _ := r.GetInt("role_id")
	o := orm.NewOrm()
	role := auth.Role{Id: role_id}
	o.QueryTable("sys_role").Filter("id", role_id).One(&role)
	// 已绑定的用户
	o.LoadRelated(&role, "User")
	// 未绑定用户
	users := []auth.User{}
	if len(role.User) > 0 {
		// 查询数据库中已经选中的用户

		o.QueryTable("sys_user").Filter("is_delete", 0).Filter("is_active", 1).Exclude("id__in", role.User).All(&users)

	} else {
		// 未绑定数据
		o.QueryTable("sys_user").Filter("is_delete", 0).Filter("is_active", 1).All(&users)

	}
	r.Data["role"] = role
	r.Data["users"] = users
	r.TplName = "auth/role-user-add.html"

}

// 角色用户配置
func (r *RoleController) DoRoleUser() {
	role_id, _ := r.GetInt("role_id")
	user_ids := r.GetString("user_ids")
	user_ids_arr := strings.Split(user_ids, ",")

	message_map := map[string]interface{}{}
	// 多对多操作 先清空 数据库已经绑定的用户,再重新添加
	role := auth.Role{Id: role_id}
	o := orm.NewOrm()
	m2m := o.QueryM2M(&role, "User")
	m2m.Clear()

	for _, user_id := range user_ids_arr {
		user := auth.User{Id: utils.StrToInt(user_id)}
		m2m := o.QueryM2M(&role, "User")
		m2m.Add(&user)

	}
	message_map["code"] = 200
	message_map["msg"] = "添加成功"
	r.Data["json"] = message_map
	r.ServeJSON()

}

func (r *RoleController) GetAuthJson() {
	role_id, _ := r.GetInt("role_id")
	o := orm.NewOrm()
	qs := o.QueryTable("sys_auth").Filter("is_delete", 0)
	// 已经绑定的权限
	role := auth.Role{Id: role_id}
	o.LoadRelated(&role, "Auth")
	auth_id_has := []int{}
	for _, auth_data := range role.Auth {
		auth_id_has = append(auth_id_has, auth_data.Id)
	}

	// 所有权限
	auths := []auth.Auth{}
	qs.All(&auths)

	auth_arr_map := []map[string]interface{}{}
	for _, auth_data := range auths {
		id := auth_data.Id
		pId := auth_data.Pid
		name := auth_data.AuthName
		if pId == 0 {
			auth_map := map[string]interface{}{"id": id, "pId": pId, "name": name, "open": false}
			auth_arr_map = append(auth_arr_map, auth_map)
		} else {
			auth_map := map[string]interface{}{"id": id, "pId": pId, "name": name}
			auth_arr_map = append(auth_arr_map, auth_map)
		}

	}
	auth_maps := map[string]interface{}{}
	auth_maps["auth_ids_has"] = auth_id_has
	auth_maps["auth_arr_map"] = auth_arr_map
	r.Data["json"] = auth_maps
	r.ServeJSON()

}
