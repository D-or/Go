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
 *     Initial: 2018/02/01      Lin Hao
 */

package taptap

import (
	"strings"
	"strconv"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"

	"../../crawler"
)

// TapTap structure
type App struct {
	Name      string
	Avatar    string
	Tag       string
	Category  []string
	Desc      string
	Grade     float64
	Developer string
	Image     string
	WebView   string
}

type appCrawler struct {
	collector *colly.Collector
	dataPipe  chan *App
	topic     *string
}

// NewAppCrawler generates a crawler for app
func NewAppCrawler(tag string, dataPipe chan *App) crawler.Crawler {
	return &appCrawler{
		collector: colly.NewCollector(),
		dataPipe:  dataPipe,
		topic:     &tag,
	}
}

// Crawler interface Init
func (c *appCrawler) Init() error {
	c.collector.OnHTML("section.app-top-list", c.parse)

	return nil
}

// Crawler interface Start
func (c *appCrawler) Start() error {
	return c.collector.Visit("https://www.taptap.com/top/" + *c.topic)
}

func (c *appCrawler) parse(e *colly.HTMLElement) {
	e.DOM.Children().Each(c.parseContent)
}

func (c *appCrawler) parseContent(_ int, s *goquery.Selection) {
	rawName := s.Children().Eq(1).Children().Eq(0).Text()
	name := strings.TrimSpace(rawName)

	rawAvatar, _ := s.Children().Eq(0).Find("img").Attr("src")
	avatar := strings.TrimSpace(rawAvatar)

	var category []string
	s.Children().Eq(1).Children().Eq(4).Children().Each(func(_ int, selection *goquery.Selection) {
		category = append(category, selection.Text())
	})

	rawDesc := s.Children().Eq(1).Children().Eq(2).Text()
	desc := strings.TrimSpace(rawDesc)

	rawGrade := s.Children().Eq(1).Children().Eq(3).Find("span").Text()
	grade, _ := strconv.ParseFloat(rawGrade, 0)

	rawDeveloper := s.Children().Eq(1).Children().Eq(1).Text()
	developer := strings.TrimSpace(rawDeveloper)

	rawImage, _ := s.Children().Eq(2).Find("img").Attr("src")
	image := strings.TrimSpace(rawImage)

	rawWebView, _ := s.Children().Eq(1).Children().Eq(0).Attr("href")
	webView := strings.TrimSpace(rawWebView)

	data := &App{
		Name:      name,
		Avatar:    avatar,
		Tag:       *c.topic,
		Category:  category,
		Desc:      desc,
		Grade:     grade,
		Developer: developer,
		Image:     image,
		WebView:   webView,
	}

	c.dataPipe <- data
}
