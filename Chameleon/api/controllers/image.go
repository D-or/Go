/*
 * MIT License
 *
 * Copyright (c) 2017 Lin Hao.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the 'Software'), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED 'AS IS', WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

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
