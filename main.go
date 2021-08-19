package main

import (
	_ "IOS_Volta/models/auth" //models文件夹下嵌套了一个文件夹 在使用命令行建表的时候需要引入包含项目的路径,否则无法迁移
	_ "IOS_Volta/models/myCenter"
	_ "IOS_Volta/routers"
	"IOS_Volta/utils"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	//数据库连接开始
	dbhost := beego.AppConfig.String("dbhost")
	dbport := beego.AppConfig.String("dbport")
	dbuser := beego.AppConfig.String("dbuser")
	dbpassword := beego.AppConfig.String("dbpassword")
	db := beego.AppConfig.String("db")
	//注册mysql Driver
	orm.RegisterDriver("mysql", orm.DRMySQL)
	//构造conn连接
	//用户名:密码@tcp(url地址)/数据库
	conn := dbuser + ":" + dbpassword + "@tcp(" + dbhost + ":" + dbport + ")/" + db + "?charset=utf8"
	//注册数据库连接
	orm.RegisterDataBase("default", "mysql", conn)
	fmt.Printf("数据库连接成功！%s\n", conn)

	//数据库连接结束
	//数据库命令行迁移
	orm.RunCommand()
	// 未登录请求拦截
	beego.InsertFilter("/main/*", beego.BeforeRouter, utils.LoginFilter)
	// 开启全局打印sql语句.
	//orm.Debug = true

	// 日志
	logs.SetLogger(logs.AdapterMultiFile, `{"filename":"logs/ions.log","separate":["error","info"]}`)
	logs.SetLogFuncCallDepth(3)
}

func main() {
	beego.Run()
}
