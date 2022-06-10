package main

import (
	"douban/douban"

	"douban/spiders"
)

func init() {
	douban.DefaultEngine.Register(spiders.DoubanSpider())
}
