/*
 * Revision History:
 *     Initial: 2018/05/23      Lin Hao
 */

package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"

	"go_trainning/Chameleon/api/models"
	_ "go_trainning/Chameleon/api/routers"
)

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)

	orm.RegisterModel(new(models.Image))

	orm.RegisterDataBase("default", "mysql", "chameleon:chameleon@/chameleon?charset=utf8")

	// if err := orm.RunSyncdb("defalut", true, true); err != nil {
	// 	panic(err)
	// }
}

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	beego.BConfig.CopyRequestBody = true

	beego.Run(":2375")
}
