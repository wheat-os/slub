package spiders

import (
	"douban/douban"
	"fmt"
	"sync"

	"github.com/antchfx/htmlquery"
	"github.com/wheat-os/slubby/spider"
	"github.com/wheat-os/slubby/stream"
)

type doubanSpider struct{}

func (t *doubanSpider) UId() string {
	return "douban"
}

func (t *doubanSpider) FQDN() string {
	return "movie.douban.com"
}

func (t *doubanSpider) Parse(response *stream.HttpResponse) (stream.Stream, error) {
	root, err := htmlquery.Parse(response.Body)
	if err != nil {
		return nil, err
	}

	lis := htmlquery.Find(root, `//ol[@class="grid_view"]/li`)

	stms := make([]stream.Stream, 0, len(lis))
	for _, li := range lis {
		name := htmlquery.Find(li, `.//div[@class="info"]//a/span/text()`)[0]

		stms = append(stms, &douban.DoubanItem{
			Item:  stream.BasicItem(t),
			Title: name.Data,
		})
	}
	return stream.StreamLists(t, stms...), nil
}

func (t *doubanSpider) StartRequest() stream.Stream {

	req, _ := stream.StreamListRangeInt(t, func(i int) (stream.Stream, error) {
		url := "https://movie.douban.com/top250?start=%d"
		return stream.Request(t, fmt.Sprintf(url, i*25), nil)
	}, 0, 9)

	return req
}

var (
	doubanOnce = sync.Once{}
	doubanP    *doubanSpider
)

func DoubanSpider() spider.Spider {
	doubanOnce.Do(func() {
		doubanP = &doubanSpider{}
	})
	return doubanP
}
