/*
 * Revision History:
 *     Initial: 2018/05/30      Lin Hao
 */

package user

import (
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

// User - The struction of user.
type User struct {
	Id      int       `json:"id"      orm:"pk;auto;index"`
	Openid  string    `json:"openid"`
	Created time.Time `json:"created" orm:"auto_now_add;type(datetime)"`
}

// Add user to the DB.
func Add(openID string) (int64, error) {
	o := orm.NewOrm()
	o.Using("default")

	u := new(User)
	u.Openid = openID

	_, id, err := o.ReadOrCreate(u, "Openid")
	if err != nil {
		beego.Error("ReadOrCreate Error: ", err)
		return -1, nil
	}

	return id, nil
}
