/*
 * Revision History:
 *     Initial: 2018/05/23      Lin Hao
 */

package image

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"

	"go_trainning/Chameleon/api/utils"
)

// Image - The struction of image.
type Image struct {
	Id      int       `json:"id"     orm:"pk;auto;index"`
	Userid  int       `json:"userID" orm:"index"`
	Name    string    `json:"name"`
	Path    string    `json:"path"`
	Created time.Time `json:"created"       orm:"auto_now_add;type(datetime)"`
}

func init() {
}

// Add a url of image to the DB.
func Add(r *http.Request) (string, int64) {
	o := orm.NewOrm()
	o.Using("default")

	fileName, err := utils.Save(r)
	if err != nil {
		beego.Error("Draw Error: ", err)
		return "", -1
	}

	var texts []string
	json.Unmarshal([]byte(r.FormValue("texts")), &texts)

	utils.Generate(fileName, texts, r.FormValue("position"), r.FormValue("color"))

	generatedPath := "https://www.doublewoodh.club/images/generated/" + fileName
	originPath := "https://www.doublewoodh.club/images/origin/" + fileName

	userID, err := strconv.Atoi(r.FormValue("userID"))
	if err != nil {
		beego.Error(err)
		return "", -1
	}

	originImage := new(Image)
	originImage.Userid = userID
	originImage.Name = fileName
	originImage.Path = originPath

	_, err = o.Insert(originImage)
	if err != nil {
		beego.Error("Insert originImage Error: ", err)
		return "", -1
	}

	generatedImage := new(Image)
	generatedImage.Userid = userID
	generatedImage.Name = fileName
	generatedImage.Path = generatedPath

	id, err := o.Insert(generatedImage)
	if err != nil {
		beego.Error("Insert generatedImage Error: ", err)
		return "", -1
	}

	return generatedPath, id
}

// // GetOne - Get the url of image by id.
// func GetOne(id string) (*Image, error) {
// 	if v, ok := Images[id]; ok {
// 		return v, nil
// 	}

// 	return nil, errors.New("ImageId Not Exist")
// }

// GetByUserID - Get all images by UserID.
func GetByUserID(userID int64) []*Image {
	var images []*Image

	o := orm.NewOrm()
	o.Using("default")

	_, err := o.QueryTable("image").Filter("Userid", userID).All(&images)
	if err != nil {
		beego.Error("Get all image Error: ", err)
		return nil
	}

	return images
}

// Update - Update the url of image by id.
func Update(id int) {
}

// Delete the url of image by id.
func Delete(id []int) error {
	o := orm.NewOrm()
	o.Using("default")
	err := o.Begin()

	qs := o.QueryTable("image")

	_, err = qs.Filter("Id", id[0]).Delete()
	if err != nil {
		err = o.Rollback()
		return err
	}

	_, err = qs.Filter("Id", id[1]).Delete()
	if err != nil {
		err = o.Rollback()
		return err
	}

	o.Commit()

	return nil
}
