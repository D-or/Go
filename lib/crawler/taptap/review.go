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
)

// TapTap structure
type Review struct {
	App       string
	AppAvatar string
	UserName  string
	Detail    string
	Vote      int
}

type reviewCrawler struct {
	collector *colly.Collector
	dataPipe  chan *Review
}

// NewReviewCrawler generates a crawler for reviews
func NewReviewCrawler(dataPipe chan *Review) crawler.Crawler {
	return &reviewCrawler{
		collector: colly.NewCollector(),
		dataPipe:  dataPipe,
	}
}

// Crawler interface Init
func (c *reviewCrawler) Init() error {
	c.collector.OnHTML("section.index-reviews>ul.list-unstyled", c.parse)

	return nil
}

// Crawler interface Start
func (c *reviewCrawler) Start() error {
	return c.collector.Visit("https://www.taptap.com")
}

func (c *reviewCrawler) parse(e *colly.HTMLElement) {
	e.DOM.Children().Each(c.parseContent)
}

func (c *reviewCrawler) parseContent(_ int, s *goquery.Selection) {
	rawApp := s.Children().Eq(0).Children().Eq(1).Find("a").Text()
	app := strings.TrimSpace(rawApp)

	rawAppAvatar, _ := s.Children().Eq(0).Find("img").Attr("src")
	appAvatar := strings.TrimSpace(rawAppAvatar)

	rawUserName := s.Children().Eq(1).Children().Eq(1).Find("a").Text()
	userName := strings.TrimSpace(rawUserName)

	rawDetail := s.Children().Eq(1).Children().Eq(0).Find("p").Text()
	detail := strings.TrimSpace(rawDetail)

	rawVote := s.Children().Eq(2).Find("span").Text()
	vote, _ := strconv.Atoi(rawVote)

	data := &Review{
		App:       app,
		AppAvatar: appAvatar,
		UserName:  userName,
		Detail:    detail,
		Vote:      vote,
	}

	c.dataPipe <- data
}
