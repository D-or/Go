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
	Id      int       `json:"id"      orm:"pk;auto;index"`
	Userid  int       `json:"userID"`
	Name    string    `json:"name"`
	Path    string    `json:"path"`
	Created time.Time `json:"created" orm:"auto_now_add;type(datetime)"`
}

// Uploaded - The struction of image uploaded.
type Uploaded struct {
	Image
}

func init() {
}

// Add a url of image to the DB.
func Add(r *http.Request) (string, int64) {
	o := orm.NewOrm()
	o.Using("default")

	fileName, err := utils.Save(r, "origin")
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

// Delete the url of image by id.
func Delete(id []int, table string) error {
	o := orm.NewOrm()
	o.Using("default")
	err := o.Begin()

	qs := o.QueryTable(table)

	for _, i := range id {
		_, err = qs.Filter("Id", i).Delete()
		if err != nil {
			err = o.Rollback()
			return err
		}
	}

	o.Commit()

	return nil
}

// GetUploaded - Get all images uploaded.
func GetUploaded() []*Uploaded {
	var uploaded []*Uploaded

	o := orm.NewOrm()
	o.Using("default")

	_, err := o.QueryTable("uploaded").All(&uploaded)
	if err != nil {
		beego.Error("Get all uploaded image Error: ", err)
		return nil
	}

	return uploaded
}

// GetUploadedByUserID - Get all images uploaded by userID.
func GetUploadedByUserID(userID int64) []*Uploaded {
	var uploaded []*Uploaded

	o := orm.NewOrm()
	o.Using("default")

	_, err := o.QueryTable("uploaded").Filter("Userid", userID).All(&uploaded)
	if err != nil {
		beego.Error("Get all uploaded image Error: ", err)
		return nil
	}

	return uploaded
}

// Upload image.
func Upload(r *http.Request, fileName string) error {
	o := orm.NewOrm()
	o.Using("default")

	userID, err := strconv.Atoi(r.FormValue("userID"))
	if err != nil {
		beego.Error(err)
		return err
	}

	uploadURL := "https://www.doublewoodh.club/images/upload/" + fileName

	uploadImage := new(Uploaded)
	uploadImage.Userid = userID
	uploadImage.Name = fileName
	uploadImage.Path = uploadURL

	_, err = o.Insert(uploadImage)
	if err != nil {
		beego.Error("Insert uploaded image Error: ", err)
		return err
	}

	return nil
}
