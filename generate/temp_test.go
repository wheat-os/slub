package generate

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetProjectModByFile(t *testing.T) {
	SetProjectModByFile("/home/bandl/go/src/github.com/wheat-os/slub/go.mod")
}

var testSpiderTemp = `package spiders

import (
	"fmt"
	"sync"

	"github.com/wheat-os/slubby/spider"
	"github.com/wheat-os/slubby/stream"
)

type firstSpider struct{}

func (t *firstSpider) UId() string {
	return "first"
}

func (t *firstSpider) FQDN() string {
	return "www.first.com"
}

func (t *firstSpider) Parse(response *stream.HttpResponse) (stream.Stream, error) {
	fmt.Println(response.Text())

	return nil, nil
}

func (t *firstSpider) StartRequest() stream.Stream {
	req, _ := stream.Request(t, "http://www.baidu.com", nil)
	return req
}

var (
	once  = sync.Once{}
	first *firstSpider
)

func FirstSpider() spider.Spider {
	once.Do(func() {
		first = &firstSpider{}
	})
	return first
}
`

func TestXpath(t *testing.T) {
	// func\s+(\w+)\s?\\(\\)\w?stream.Stream\w?{
	reg, err := regexp.Compile(`func\s+(\w+)\w?\s?\(\)\s?spider.Spider`)
	require.NoError(t, err)

	fmt.Println(reg.FindAllStringSubmatch(testSpiderTemp, 1))
}
