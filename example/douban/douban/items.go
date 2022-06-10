package douban

import "github.com/wheat-os/slubby/stream"

type DoubanItem struct {
	stream.Item
	Title string `csv:"title"`
}
