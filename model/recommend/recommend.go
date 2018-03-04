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

package recommend

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
		cname = "recommend"
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

	go Service.Save()
}

type Recommend struct {
	ID            bson.ObjectId `bson:"_id,omitempty"`
	Name          string        `bson:"name"`
	Avatar        string        `bson:"avatar"`
	Image         string        `bson:"image"`
	Desc          string        `bson:"desc"`
	Grade         float64       `bson:"grade"`
	Comments      int           `bson:"comments"`
}

// Save - save the data of recommends to db
func (sp *serviceProvider) Save() error {
	conn := session.Connect()
	defer conn.Disconnect()

	dataPipe := make(chan *taptap.Recommend)
	c := taptap.NewRecommendCrawler(dataPipe)

	go func() {
		_, err := conn.DeleteAll(nil)
		if err != nil {
			return
		}

		for {
			data := <-dataPipe

			recommendInfo := &Recommend{
				Name:          data.Name,
				Avatar:        data.Avatar,
				Image:         data.Image,
				Desc:          data.Desc,
				Grade:         data.Grade,
				Comments:      data.Comments,
			}

			err := conn.Insert(&recommendInfo)
			if err != nil {
				return
			}
		}
	}()

	crawler.StartCrawler(c)

	return nil
}

// List displays all recommends
func (sp *serviceProvider) Get() ([]Recommend, error) {
	var (
		err error
		result []Recommend
	)

	conn := session.Connect()
	defer conn.Disconnect()

	err = conn.GetMany(nil, &result)

	return result, err
}
