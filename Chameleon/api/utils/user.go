/*
 * Revision History:
 *     Initial: 2018/05/31      Lin Hao
 */

package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/astaxie/beego"
)

// GetOpenID - Get User OpenID by code.
func GetOpenID(code string) (string, error) {
	var response map[string]interface{}

	res, err := http.Get("https://api.weixin.qq.com/sns/jscode2session?appid=wx605779143a485a34&secret=d0a0063540f9a47d8601eb2d95ae111e&js_code=" + code + "&grant_type=authorization_code")
	if err != nil {
		beego.Debug("Get OpenID Error: ", err)
		return "", err
	}
	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		beego.Debug("Read Body Error: ", err)
		return "", err
	}

	json.Unmarshal(resBody, &response)

	if openID, ok := response["openid"].(string); ok {
		return openID, nil
	}

	return "", err
}
