/*
 * Revision History:
 *     Initial: 2018/05/23      Lin Hao
 */

package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"

	"go_trainning/Chameleon/api/models/image"
	"go_trainning/Chameleon/api/models/user"
	_ "go_trainning/Chameleon/api/routers"
)

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)

	orm.RegisterModel(new(image.Image))
	orm.RegisterModel(new(user.User))

	orm.RegisterDataBase("default", "mysql", "chameleon:chameleon@/chameleon?charset=utf8")
}

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	beego.BConfig.CopyRequestBody = true

	beego.Run(":2375")
}
