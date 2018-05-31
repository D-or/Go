/*
 * Revision History:
 *     Initial: 2018/05/30      Lin Hao
 */

package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego"

	"go_trainning/Chameleon/api/models/user"
	"go_trainning/Chameleon/api/utils"
)

// UserController - Operations about user.
type UserController struct {
	beego.Controller
}

// Login with code from frontend.
func (uc *UserController) Login() {
	var (
		req    map[string]interface{}
		userID int64
	)

	json.Unmarshal(uc.Ctx.Input.RequestBody, &req)
	code := req["code"].(string)

	openID, err := utils.GetOpenID(code)
	if err != nil {
		beego.Debug("GetOpenID Error: ", err)

		uc.Data["json"] = map[string]interface{}{
			"userID": -1,
		}

		goto finish
	}

	userID, err = user.Add(openID)

	uc.Data["json"] = map[string]interface{}{
		"userID": userID,
	}

finish:
	uc.ServeJSON()
}
