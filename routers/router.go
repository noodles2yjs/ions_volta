package routers

import (
	"IOS_Volta/controllers"
	"IOS_Volta/controllers/auth"
	"IOS_Volta/controllers/cars"
	"IOS_Volta/controllers/login"
	"IOS_Volta/controllers/user"
	"github.com/astaxie/beego"
)

func init() {
	//beego.Router("/", &controllers.MainController{})
	beego.Router("/", &login.LoginController{})
	beego.Router("/user/log_out", &login.LoginController{}, "get:LogOut")
	beego.Router("/change_captcha", &login.LoginController{}, "get:ChangeCaptcha")
	// 需要登陆 后台router

	// 后台首页
	beego.Router("/main/index", &controllers.HomeController{})
	// 消息通知 可以用websocket 发送消息,触发使用定时触发 这里不会使用,使用登陆后自动进入对的方法改变
	beego.Router("/main/index/notify/list", &controllers.HomeController{}, "get:NotifyList")
	beego.Router("/main/index/notify/readNotify", &controllers.HomeController{}, "get:ReadNotify")

	beego.Router("/main/welcome", &controllers.HomeController{}, "get:Welcome")

	// user模块
	beego.Router("/main/user/is_active", &user.UserController{}, "post:IsActive")
	beego.Router("/main/user/list", &user.UserController{}, "get:List")
	beego.Router("/main/user/to_add", &user.UserController{}, "get:ToAdd")
	beego.Router("/main/user/do_add", &user.UserController{}, "post:DoAdd")
	beego.Router("/main/user/to_update", &user.UserController{}, "get:ToUpdate")
	beego.Router("/main/user/do_update", &user.UserController{}, "post:DoUpdate")
	beego.Router("/main/user/delete", &user.UserController{}, "post:Delete")
	beego.Router("/main/user/reset_password", &user.UserController{}, "post:ResetPassword")
	beego.Router("/main/user/multi_delete", &user.UserController{}, "post:MultiDelete")

	//权限模块
	beego.Router("/main/auth/list", &auth.AuthController{}, "get:List")
	beego.Router("/main/auth/is_active", &auth.AuthController{}, "post:IsActive")
	beego.Router("/main/auth/to_add", &auth.AuthController{}, "get:ToAdd")
	beego.Router("/main/auth/do_add", &auth.AuthController{}, "post:DoAdd")
	beego.Router("/main/auth/to_update", &auth.AuthController{}, "get:ToUpdate")
	beego.Router("/main/auth/do_update", &auth.AuthController{}, "post:DoUpdate")
	beego.Router("/main/auth/delete", &auth.AuthController{}, "post:Delete")
	beego.Router("/main/auth/multi_delete", &auth.AuthController{}, "post:MultiDelete")
	// 角色模块
	beego.Router("/main/role/list", &auth.RoleController{}, "get:List")
	beego.Router("/main/role/to_add", &auth.RoleController{}, "get:ToAdd")
	beego.Router("/main/role/do_add", &auth.RoleController{}, "post:DoAdd")
	// 角色-权限配置
	beego.Router("/main/role/to_roleAuth", &auth.RoleController{}, "get:ToRoleAuth")
	beego.Router("/main/role/do_roleAuth", &auth.RoleController{}, "post:DoRoleAuth")
	beego.Router("/main/role/get_authJson", &auth.RoleController{}, "get:GetAuthJson")
	// 角色-用户配置
	beego.Router("/main/role/to_roleUser", &auth.RoleController{}, "get:ToRoleUser")
	beego.Router("/main/role/do_roleUser", &auth.RoleController{}, "post:DoRoleUser")

	// 个人中心
	beego.Router("/main/user/my_center", &user.MyCenterController{})
	beego.Router("/main/user/salary_slip", &user.SalarySlipController{})
	beego.Router("/main/user/salary_slip_detail", &user.SalarySlipController{}, "get:Detail")

	// 车辆管理模块
	// 车辆品牌管理
	beego.Router("/main/carBrand/list", &cars.CarBrandController{}, "get:List")
	beego.Router("/main/carBrand/toAdd", &cars.CarBrandController{}, "get:ToAdd")
	beego.Router("/main/carBrand/doAdd", &cars.CarBrandController{}, "post:DoAdd")
	beego.Router("/main/carBrand/delete", &cars.CarBrandController{}, "post:Delete")
	beego.Router("/main/carBrand/toUpdate", &cars.CarBrandController{}, "get:ToUpdate")
	beego.Router("/main/carBrand/doUpdate", &cars.CarBrandController{}, "post:DoUpdate")
	beego.Router("/main/carBrand/IsActive", &cars.CarBrandController{}, "post:IsActive")
	beego.Router("/main/carBrand/muli_delete", &cars.CarBrandController{}, "post:MuliDelete")

	// 车辆管理
	beego.Router("/main/car/list", &cars.CarController{}, "get:List")
	beego.Router("/main/car/ToAdd", &cars.CarController{}, "get:ToAdd")
	beego.Router("/main/car/DoAdd", &cars.CarController{}, "post:DoAdd")
	beego.Router("/main/car/ToUpdate", &cars.CarController{}, "get:ToUpdate")
	beego.Router("/main/car/DoUpdate", &cars.CarController{}, "post:DoUpdate")
	beego.Router("/main/car/IsActive", &cars.CarController{}, "post:IsActive")
	beego.Router("/main/car/delete", &cars.CarController{}, "post:Delete")
	beego.Router("/main/car/muli_delete", &cars.CarController{}, "post:MuliDelete")

	// 车辆申请

	beego.Router("/main/car_apply/List", &cars.CarApplyController{}, "get:List")
	beego.Router("/main/car_apply/toApply", &cars.CarApplyController{}, "get:ToApply")
	beego.Router("/main/car_apply/doApply", &cars.CarApplyController{}, "post:DoApply")
	beego.Router("/main/car_apply/myApply", &cars.CarApplyController{}, "get:MyApply")
	beego.Router("/main/car_audit/auditApply", &cars.CarApplyController{}, "get:AuditApply")
	beego.Router("/main/car_audit/toAuditApply", &cars.CarApplyController{}, "get:ToAuditApply")
	beego.Router("/main/car_audit/doAuditApply", &cars.CarApplyController{}, "post:DoAuditApply")
	beego.Router("/main/car/doReturn", &cars.CarApplyController{}, "get:DoReturn")

}
