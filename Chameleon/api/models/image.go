/*
 * Revision History:
 *     Initial: 2018/05/23      Lin Hao
 */

package models

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"

	"go_trainning/Chameleon/api/utils"
)

// Image - The struction of image.
type Image struct {
	ID      int       `json:"id" orm:"pk;auto;index"`
	Name    string    `json:"name"`
	Path    string    `json:"path"`
	Created time.Time `orm:"auto_now_add;type(datetime)"`
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

	utils.Generate(fileName, texts)

	generatedPath := "https://www.doublewoodh.club/images/generated/" + fileName
	originPath := "https://www.doublewoodh.club/images/origin/" + fileName

	originImage := new(Image)
	originImage.Name = fileName
	originImage.Path = originPath

	_, err = o.Insert(originImage)
	if err != nil {
		beego.Error("Insert originImage Error: ", err)
		return "", -1
	}

	generatedImage := new(Image)
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

// GetAll - Get All images.
func GetAll() []*Image {
	var images []*Image

	o := orm.NewOrm()
	o.Using("default")

	_, err := o.QueryTable("image").All(&images)
	if err != nil {
		return nil
	}

	return images
}

// Update - Update the url of image by id.
func Update(id string) {
}

// Delete the url of image by id.
func Delete(id string) {
}
