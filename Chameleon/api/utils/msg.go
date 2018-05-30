/*
 * Revision History:
 *     Initial: 2018/05/30      Lin Hao
 */

package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

func MsgSecCheck(texts string) error {
	var (
		getResult  map[string]interface{}
		postResult map[string]interface{}
	)

	getResp, err := http.Get("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=wx605779143a485a34&secret=d0a0063540f9a47d8601eb2d95ae111e")
	if err != nil {
		return err
	}
	defer getResp.Body.Close()

	getRespBody, err := ioutil.ReadAll(getResp.Body)
	if err != nil {
		return err
	}

	json.Unmarshal(getRespBody, &getResult)
	accessToken := getResult["access_token"].(string)

	b, err := json.Marshal(map[string]string{
		"content": texts,
	})
	body := bytes.NewBuffer([]byte(b))

	postResp, err := http.Post(
		"https://api.weixin.qq.com/wxa/msg_sec_check?access_token="+accessToken,
		"application/json",
		body,
	)
	if err != nil {
		return err
	}
	defer postResp.Body.Close()

	postRespBody, err := ioutil.ReadAll(postResp.Body)
	if err != nil {
		return err
	}

	json.Unmarshal(postRespBody, &postResult)

	if postResult["errcode"] != 0 {
		return errors.New("Content is risky")
	}

	return nil
}
