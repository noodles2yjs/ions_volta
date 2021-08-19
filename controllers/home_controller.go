package controllers

import (
	"IOS_Volta/models/auth"
	"IOS_Volta/utils"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type HomeController struct {
	beego.Controller
}

func (h *HomeController) Get() {
	// 后端首页
	o := orm.NewOrm()

	user_id := h.GetSession("id")

	user := auth.User{Id: user_id.(int)}

	// 通过user_id  查询绑定在用户上的角色,在通过角色,查找权限,最终找出该用户所有的权限
	o.LoadRelated(&user, "Role")
	auth_arr := []int{}

	for _, role := range user.Role {
		role_data := auth.Role{Id: role.Id}
		o.LoadRelated(&role_data, "Auth")
		for _, auth_data := range role_data.Auth {
			auth_arr = append(auth_arr, auth_data.Id)
		}
	}

	qs := o.QueryTable(new(auth.Auth))
	auths := []auth.Auth{}
	trees := []auth.Tree{}
	// 获取当前用户下的第一级菜单
	qs.Filter("pid", 0).Filter("id__in", auth_arr).OrderBy("-Weight").All(&auths)

	for _, auth_data := range auths {
		pid := auth_data.Id
		tree_data := auth.Tree{Id: auth_data.Id, AuthName: auth_data.AuthName, UrlFor: auth_data.UrlFor, Weight: auth_data.Weight, Children: []*auth.Tree{}}
		GetChildNode(pid, &tree_data)
		trees = append(trees, tree_data)
	}
	o.QueryTable(new(auth.User)).Filter("id", user_id).One(&user)
	h.Data["user"] = user
	h.Data["trees"] = trees

	h.TplName = "index.html"
}
func (h *HomeController) NotifyList() {

	// 需要分页

	current_page, _ := h.GetInt("page", 1)
	// 查询
	kw := h.GetString("kw")
	pagesize := 2

	// 这是分页封装的方法，我会放在最下面(这个是重点) agricultures--对应要查询的models
	count, agricultures := NotifyListPageModel(current_page, pagesize, kw)

	// 这是封装的方法，我会放在最下面
	num_pages := utils.GetPageNum(count, pagesize)

	// 这是封装的方法，我会放在最下面
	left_pages, right_pages, left_has_more, right_has_more := utils.Get_pagination_data(num_pages, current_page, 2)
	has_previous, has_next, previous_page_number, next_page_number := utils.HasNext(current_page, num_pages)

	h.Data["left_pages"] = left_pages         //左边页码数[]int
	h.Data["right_pages"] = right_pages       //右边页码数[]int
	h.Data["left_has_more"] = left_has_more   //左边是否有页码 bool
	h.Data["right_has_more"] = right_has_more //右边是否有页码 bool

	h.Data["has_previous"] = has_previous                 //是否有上一页 bool
	h.Data["has_next"] = has_next                         //是否有下一页 bool
	h.Data["previous_page_number"] = previous_page_number //上一页页码 int
	h.Data["next_page_number"] = next_page_number         //下一页页码 int

	h.Data["num_pages"] = num_pages //总页数
	h.Data["page"] = current_page   //当前页

	h.Data["count"] = count //总数量
	h.Data["nofities"] = agricultures
	h.Data["kw"] = kw
	h.TplName = "notify_list.html"

}
func (h *HomeController) ReadNotify() {
	id, _ := h.GetInt("id")
	o := orm.NewOrm()
	qs := o.QueryTable(new(auth.MessageNotify))
	qs.Filter("id", id).Update(orm.Params{
		"ReadTag": 1,
	})

	message_data := auth.MessageNotify{}
	qs.Filter("id", id).One(&message_data)
	h.Data["message_data"] = message_data
	h.TplName = "notify_detail.html"

}

func NotifyListPageModel(page, pageSize int, kw string) (count int64, agricultures []auth.MessageNotify) {
	o := orm.NewOrm()

	qs := o.QueryTable(new(auth.MessageNotify))

	if kw != "" { // 有查询条件的
		// 总数
		count, _ = qs.Filter("IsDelete", 0).Filter("Name__contains", kw).Count()
	} else {
		qs.Filter("Name__contains", kw).Limit(pageSize).Offset((page - 1) * pageSize).RelatedSel().All(&agricultures)
		count, _ = qs.Count()
		qs.Limit(pageSize).Offset((page - 1) * pageSize).RelatedSel().All(&agricultures)

	}

	return count, agricultures
}
func (h *HomeController) Welcome() {
	h.TplName = "welcome.html"
}

// 递归 获得子菜单
func GetChildNode(pid int, treenode *auth.Tree) {
	o := orm.NewOrm()
	qs := o.QueryTable("sys_auth")
	auths := []auth.Auth{}
	_, err := qs.Filter("pid", pid).OrderBy("-weight").All(&auths)
	if err != nil {
		return
	}

	// 查询三级及以上菜单
	//for i:= 0; i<len(auths);i++{
	//	pid := auths[i].Id   // 根据pid获取所有的子解点
	//	tree_data := auth.Tree{Id:auths[i].Id,AuthName:auths[i].AuthName,UrlFor:auths[i].UrlFor,Weight:auths[i].Weight,Children:[]*auth.Tree{}}
	//	treenode.Children = append(treenode.Children,&tree_data)
	//	GetChildNode(pid,&tree_data)
	//}
	// 查询三级及以上菜单
	for _, v := range auths {
		pid := v.Id
		tree_data := auth.Tree{
			Id:       v.Id,
			AuthName: v.AuthName,
			UrlFor:   v.UrlFor,
			Weight:   v.Weight,
			Children: []*auth.Tree{},
		}
		treenode.Children = append(treenode.Children, &tree_data)

		GetChildNode(pid, &tree_data)
	}

	return
}
