package douban

import (
	"github.com/wheat-os/slubby/stream"
	"github.com/wheat-os/wlog"
)

type DoubanPipline struct{}

func (t *DoubanPipline) OpenSpider() error {
	return nil
}

func (t *DoubanPipline) CloseSpider() error {
	return nil
}

func (t *DoubanPipline) ProcessItem(item stream.Item) stream.Item {
	wlog.Info(item)
	return item
}
