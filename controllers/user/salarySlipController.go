package user

import (
	"IOS_Volta/models/auth"
	"IOS_Volta/models/myCenter"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type SalarySlipController struct {
	beego.Controller
}

func (s *SalarySlipController) Get() {
	month := s.GetString("month")
	if month == "" {
		month = time.Now().Format("2006-01")
	}
	userId := s.GetSession("id")

	o := orm.NewOrm()
	user := auth.User{}
	o.QueryTable(new(auth.User)).Filter("id", userId).One(&user)
	cardId := user.CardId
	salarySlip := myCenter.SalarySlip{}
	o.QueryTable(new(myCenter.SalarySlip)).Filter("CardId", cardId).Filter("PayDate", month).One(&salarySlip)
	s.Data["salarySlip"] = salarySlip
	s.TplName = "user/salary_slip_list.html"
}
func (s *SalarySlipController) Detail() {
	id, _ := s.GetInt("id")
	o := orm.NewOrm()
	salarySlip := myCenter.SalarySlip{}
	o.QueryTable(new(myCenter.SalarySlip)).Filter("id", id).One(&salarySlip)
	s.Data["salary_slip"] = salarySlip
	s.TplName = "user/salary_slip_detail.html"

}
