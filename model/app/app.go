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
 *     Initial: 2018/02/03      Lin Hao
 */

package app

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"../../lib/crawler/taptap"
	"../../lib/crawler"
	"github.com/fengyfei/gu/libs/mongo"
)

type serviceProvider struct{}

var (
	Service *serviceProvider
	session *mongo.Connection
)

func init() {
	const (
		cname = "app"
	)

	s, err := mgo.Dial("mongodb://127.0.0.1:27017/mosasaur")
	if err != nil {
		panic(err)
	}

	s.SetMode(mgo.Monotonic, true)
	s.DB("mosasaur").C(cname).EnsureIndex(mgo.Index{
		Key:        []string{"Name", "Developer"},
		Unique:     true,
		Background: true,
		Sparse:     true,
	})

	session = mongo.NewConnection(s, "mosasaur", cname)
	Service = &serviceProvider{}

	tags := []string{"download", "ios", "new", "played"}
	for _, tag := range tags {
		go Service.Save(tag)
	}
}

type App struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	Name      string        `bson:"name"`
	Avatar    string        `bson:"avatar"`
	Tag       string        `bson:"tag"`
	Category  []string      `bson:"category"`
	Desc      string        `bson:"desc"`
	Grade     float64       `bson:"grade"`
	Developer string        `bson:"developer"`
	Image     string        `bson:"image"`
	WebView   string        `bson:"webView"`
}

// Save - save the data of apps to db
func (sp *serviceProvider) Save(tag string) error {
	conn := session.Connect()
	defer conn.Disconnect()

	dataPipe := make(chan *taptap.App)
	c := taptap.NewAppCrawler(tag, dataPipe)

	go func() {
		_, err := conn.DeleteAll(nil)
		if err != nil {
			return
		}

		for {
			data := <-dataPipe

			if len(data.Category) == 0 {
				data.Category = []string{"无"}
			}

			appInfo := &App{
				Name:      data.Name,
				Avatar:    data.Avatar,
				Tag:       tag,
				Category:  data.Category,
				Desc:      data.Desc,
				Grade:     data.Grade,
				Developer: data.Developer,
				Image:     data.Image,
				WebView:   data.WebView,
			}

			err := conn.Insert(&appInfo)
			if err != nil {
				return
			}
		}
	}()

	crawler.StartCrawler(c)

	return nil
}

// List displays all apps
func (sp *serviceProvider) GetByTag(tag string) ([]App, error) {
	var (
		err error
		result []App
	)

	conn := session.Connect()
	defer conn.Disconnect()

	err = conn.GetMany(bson.M{"tag": tag}, &result, "-grade")

	return result, err
}
