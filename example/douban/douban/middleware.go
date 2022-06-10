package douban

import (
	"github.com/wheat-os/slubby/download/middle"
	"github.com/wheat-os/slubby/stream"
)

type DoubanMiddle struct{}

func (s *DoubanMiddle) BeforeDownload(m *middle.M, req *stream.HttpRequest) (*stream.HttpRequest, error) {
	return req, nil
}

func (s *DoubanMiddle) AfterDownload(
	m *middle.M,
	req *stream.HttpRequest,
	resp *stream.HttpResponse,
) (*stream.HttpResponse, error) {
	return resp, nil
}

func (s *DoubanMiddle) ProcessErr(m *middle.M, req *stream.HttpRequest, err error) {}
