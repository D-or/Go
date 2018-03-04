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
 *     Initial: 2018/02/04      Lin Hao
 */

package taptap

import (
	"strings"
	"strconv"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"

	"go_trainning/echo/mosasaur/lib/crawler"
)

// TapTap structure
type Recommend struct {
	Name          string
	Avatar        string
	Image         string
	Desc          string
	Grade         float64
	Comments      int
}

type recommendCrawler struct {
	collector *colly.Collector
	dataPipe  chan *Recommend
}

// NewRecommendCrawler generates a crawler for recommend
func NewRecommendCrawler(dataPipe chan *Recommend) crawler.Crawler {
	return &recommendCrawler{
		collector: colly.NewCollector(),
		dataPipe:  dataPipe,
	}
}

// Crawler interface Init
func (c *recommendCrawler) Init() error {
	c.collector.OnHTML("#recList", c.parse)

	return nil
}

// Crawler interface Start
func (c *recommendCrawler) Start() error {
	return c.collector.Visit("https://www.taptap.com")
}

func (c *recommendCrawler) parse(e *colly.HTMLElement) {
	e.DOM.Children().Filter("div.feed-rec.collapse.in").Each(c.parseContent)
}

func (c *recommendCrawler) parseContent(_ int, s *goquery.Selection) {
	rawName := s.Children().Eq(1).Text()
	name := strings.TrimSpace(rawName)

	rawAvatar, _ := s.Children().Eq(0).Find("img").Attr("src")
	avatar := strings.TrimSpace(rawAvatar)

	rawDesc := s.Children().Eq(3).Text()
	desc := strings.TrimSpace(rawDesc)

	rawGrade := s.Children().Eq(2).Children().Eq(1).Children().Eq(1).Text()
	grade, _ := strconv.ParseFloat(rawGrade, 0)

	rawImage, _ := s.Children().Eq(2).Find("img").Attr("data-src")
	image := strings.TrimSpace(rawImage)

	rawComments := s.Children().Eq(2).Children().Eq(1).Children().Eq(3).Text()
	comments, _ := strconv.Atoi(rawComments)

	data := &Recommend{
		Name:          name,
		Avatar:        avatar,
		Image:         image,
		Desc:          desc,
		Grade:         grade,
		Comments:      comments,
	}

	c.dataPipe <- data
}
