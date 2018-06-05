/*
 * Revision History:
 *     Initial: 2018/05/23      Lin Hao
 */

package routers

import (
	"github.com/astaxie/beego"

	"go_trainning/Chameleon/api/controllers"
)

func init() {
	// user api
	beego.Router("/user/login", &controllers.UserController{}, "post:Login")

	// image api
	beego.Router("/image/getall", &controllers.ImageController{}, "post:GetAll")
	beego.Router("/image/generate", &controllers.ImageController{}, "post:Generate")
	beego.Router("/image/delete", &controllers.ImageController{}, "post:Delete")
	beego.Router("/image/upload", &controllers.ImageController{}, "post:Upload")

	beego.Router("/image/uploaded/getall", &controllers.ImageController{}, "get:GetUploaded")
}
