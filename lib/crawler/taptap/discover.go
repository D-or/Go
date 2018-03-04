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
 *     Initial: 2018/02/05      Lin Hao
 */

package taptap

import (
	"strings"
	"strconv"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"

	"go_trainning/echo/mosasaur/lib/crawler"
	"gopkg.in/mgo.v2/bson"
)

// Discover structure
type Discover struct {
	Tag      string
	List     []Item
}

// TagList - the list of tag structure
type Item struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	Name     string
	Avatar   string
	Category string
	Grade    float64
}

type discoverCrawler struct {
	collector *colly.Collector
	dataPipe  chan *Discover
}

// NewDiscoverCrawler generates a crawler for discover
func NewDiscoverCrawler(dataPipe chan *Discover) crawler.Crawler {
	return &discoverCrawler{
		collector: colly.NewCollector(),
		dataPipe:  dataPipe,
	}
}

// Crawler interface Init
func (c *discoverCrawler) Init() error {
	c.collector.OnHTML("section.app-categories-simple", c.parse)

	return nil
}

// Crawler interface Start
func (c *discoverCrawler) Start() error {
	return c.collector.Visit("https://www.taptap.com/categories")
}

func (c *discoverCrawler) parse(e *colly.HTMLElement) {
	e.DOM.Each(c.parseContent)
}

func (c *discoverCrawler) parseContent(_ int, s *goquery.Selection) {
	var (
		list []Item
	)

	rawTag := s.Children().Eq(0).Children().Eq(0).Text()
	tag := strings.TrimSpace(rawTag)

	s.Children().Eq(1).Children().Each(func(_ int, selection *goquery.Selection) {
		rawName := selection.Children().Eq(1).Children().Eq(0).Text()
		name := strings.TrimSpace(rawName)

		rawAvatar, _ := selection.Children().Eq(0).Find("img").Attr("data-src")
		avatar := strings.TrimSpace(rawAvatar)

		rawCategory := selection.Children().Eq(1).Children().Eq(1).Children().Eq(0).Text()
		category := strings.TrimSpace(rawCategory)

		rawGrade := selection.Children().Eq(1).Children().Eq(1).Children().Eq(1).Children().Eq(0).Text()
		grade, _ := strconv.ParseFloat(rawGrade, 0)

		item := &Item{
			Name:     name,
			Avatar:   avatar,
			Category: category,
			Grade:    grade,
		}

		list = append(list, *item)
	})

	data := &Discover{
		Tag:  tag,
		List: list,
	}

	c.dataPipe <- data
}
