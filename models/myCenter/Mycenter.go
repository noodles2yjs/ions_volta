package myCenter

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type SalarySlip struct {
	Id                int       `orm :"pk;auto"`
	CardId            string    `orm:"size(64);description(员工工号);"`
	Basepay           float64   `orm:"description(基本工资);digits(12);decimals(2)"`
	WorkingDays       float64   `orm:"description( 工作天数);digits(3);decimals(1)"`
	DaysOff           float64   `orm:"description(请假天数);digits(3);decimals(1)"`
	DaysOffNo         float64   `orm:"description(调休天数);digits(3);decimals(1);column(daysOffNo)"`
	Reward            float64   `orm:"description(奖金);digits(8);default(0);decimals(2)"`
	RentSubSidy       float64   `orm:"description(租房补贴);default(0);digits(6);decimals(2);column(rentSubSidy)"`
	TransSubSidy      float64   `orm:"description(交通补贴);digits(6);default(0);decimals(2);column(transSubSidy)"`
	SocialSecurity    float64   `orm:"description(社保);digits(6);default(0);decimals(2);column(socialSecurity)"`
	HousProvidentFund float64   `orm:"description(住房公积金);digits(6);decimals(2);default(0);column(housProvidentFund)"`
	PersonalPncomeTax float64   `orm:"description(个税);digits(6);decimals(2);column(personalPncomeTax)"`
	Fine              float64   `orm:"description(罚金);digits(6);decimals(2);default(0);"`
	NetSalary         float64   `orm:"description(实发工资);digits(10);decimals(2);default(0);column(netSalary)"`
	PayDate           string    `orm:"description(工资月份);size(32);column(payDate)"`
	CreateTime        time.Time `orm:"type(datetime);auto_now;description(创建时间);column(createTime)"`
}

func (s *SalarySlip) TableName() string {
	return "salarySlip"
}

func init() {
	orm.RegisterModel(new(SalarySlip))

}
