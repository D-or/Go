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
	beego.Router("/image/getall", &controllers.ImageController{}, "get:GetAll")
	beego.Router("/image/generate", &controllers.ImageController{}, "post:Generate")
}
