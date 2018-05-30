/*
 * Revision History:
 *     Initial: 2018/05/23      Lin Hao
 */

package controllers

import (
	"github.com/astaxie/beego"

	"go_trainning/Chameleon/api/models"
	"go_trainning/Chameleon/api/utils"
)

// ImageController - Operations about image.
type ImageController struct {
	beego.Controller
}

// GetAll - Get all images.
func (ic *ImageController) GetAll() {
	images := models.GetAll()

	ic.Data["json"] = map[string]interface{}{
		"images": images,
	}

	ic.ServeJSON()
}

// Generate the image.
func (ic *ImageController) Generate() {
	var (
		path string
		id   int64
	)

	err := utils.MsgSecCheck(ic.Ctx.Request.FormValue("texts"))
	if err != nil {
		ic.Data["json"] = map[string]interface{}{
			"imageId": 0,
			"image":   "",
		}

		goto finish
	}

	path, id = models.Add(ic.Ctx.Request)
	if id == -1 {
		ic.Data["json"] = map[string]interface{}{
			"imageId": -1,
			"image":   "",
		}

		goto finish
	}

	ic.Data["json"] = map[string]interface{}{
		"imageId": id,
		"image":   path,
	}

finish:
	ic.ServeJSON()
}
