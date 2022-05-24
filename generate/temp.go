package generate

var (
	projectPath string
	projectMod  string
)

func SetProjectPath(p string) {
	projectPath = p
}

func SetProjectMod(p string) {
	projectMod = p
}

var spiderTemp = `package spiders

import (
	"fmt"
	"sync"

	"github.com/wheat-os/slubby/spider"
	"github.com/wheat-os/slubby/stream"
)

type {{.nameCap}}Spider struct{}

func (t *{{.nameCap}}Spider) UId() string {
	return "{{.uid}}"
}

func (t *{{.nameCap}}Spider) FQDN() string {
	return "{{.fqdn}}"
}

func (t *{{.nameCap}}Spider) Parse(response *stream.HttpResponse) (stream.Stream, error) {
	fmt.Println(response.Text())

	return nil, nil
}

func (t *{{.nameCap}}Spider) StartRequest() stream.Stream {
	req, _ := stream.Request(t, "http://www.baidu.com", nil)
	return req
}

var (
	once = sync.Once{}
	{{.nameCap}} *{{.nameCap}}Spider
)

func New{{.nameCap}}Spider() spider.Spider {
	once.Do(func() {
		{{.nameCap}} = &{{.nameCap}}Spider{}
	})
	return {{.nameCap}}
}
`

var settingTemp = `package temp

import (
	"github.com/wheat-os/slubby/download"
	"github.com/wheat-os/slubby/download/middle"
	"github.com/wheat-os/slubby/engine"
	"github.com/wheat-os/slubby/outputter"
	"github.com/wheat-os/slubby/outputter/pipline"
	"github.com/wheat-os/slubby/scheduler"
	"github.com/wheat-os/slubby/scheduler/buffer"
	"github.com/wheat-os/slubby/scheduler/filter"
	"github.com/wheat-os/wlog"
)

// ***************************************** Logger *****************************************
func init() {
	wlog.SetStdOptions(wlog.WithDisPlayLevel(wlog.InfoLevel))
	wlog.SetStdOptions(wlog.WithDisableCaller(true))
}

// **************************************** Scheduler ***************************************

// filter
var filterModule = filter.ShortBloomFilter()

// buffer
var bufferModule = buffer.ShortQueue()

var schedulerModule = scheduler.ShortScheduler(
	scheduler.WithFilter(filterModule),
	scheduler.WithBuffer(bufferModule),
)

// **************************************** Download ***************************************
// download middle
var downloadMiddlewareModule = middle.MiddleGroup(
	middle.LogMiddle(),
)

// download limiter
// var downloadLimiterModule = limiter.ShortLimiter(
// 	time.Second * 3,
// )

var downloadModule = download.ShortDownload(
	// download.WithLimiter(downloadLimiterModule),
	download.WithDownloadMiddle(downloadMiddlewareModule),
)

// *************************************** Outputter **************************************

// pipline
var piplineModule = pipline.GroupPipline(
	&{{.nameCap}}Pipline{},
)

var outputterModule = outputter.ShortOutputter(
	outputter.WithPipline(piplineModule),
)

// ***************************************** Engine ***************************************
var DefaultEngine = engine.ShortEngine(
	engine.WithScheduler(schedulerModule),
	engine.WithDownload(downloadModule),
	engine.WithOutputter(outputterModule),
)
`

var itemTemp = `package {{.pkgName}}

import "github.com/wheat-os/slubby/stream"

type {{.nameCap}}Item struct {
	stream.Item
	Desc string
}
`

var piplineTemp = `package {{.pkgName}}

import (
	"github.com/wheat-os/slubby/stream"
)

type {{.nameCap}}Pipline struct{}

func (t *{{.nameCap}}Pipline) OpenSpider() error {
	return nil
}

func (t *{{.nameCap}}Pipline) CloseSpider() error {
	return nil
}

func (t *{{.nameCap}}Pipline) ProcessItem(item stream.Item) stream.Item {
	return item
}
`

var middlewareTemp = `package {{.pkgName}}

import (
	"github.com/wheat-os/slubby/download/middle"
	"github.com/wheat-os/slubby/stream"
)

type {{.nameCap}}Middle struct{}

func (s *{{.nameCap}}Middle) BeforeDownload(m *middle.M, req *stream.HttpRequest) (*stream.HttpRequest, error) {
	return req, nil
}

func (s *{{.nameCap}}Middle) AfterDownload(
	m *middle.M,
	req *stream.HttpRequest,
	resp *stream.HttpResponse,
) (*stream.HttpResponse, error) {
	return resp, nil
}

func (s *{{.nameCap}}Middle) ProcessErr(m *middle.M, req *stream.HttpRequest, err error) {}
`

var confTem = `[slub]
version = "1.0.0"

[slubby]
version = "1.0.0"
`

var mainTemp = `package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"{{.mod}}/{{.pkgName}}"
	"github.com/wheat-os/wlog"
	"github.com/spf13/viper"
)

var projectConfFile = "{{.confFile}}"

func signalClose(cannel context.CancelFunc) {

	sig := make(chan os.Signal, 1)
	// 监听 退出
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	<-sig

	wlog.Info("the listener hears the exit signal, exiting, please do not kill the process directly.")
	cannel()
}

func main() {
	viper.SetConfigFile(projectConfFile)

	engine := {{.pkgName}}.DefaultEngine
	ctx, cannel := context.WithCancel(context.Background())
	go signalClose(cannel)

	engine.Start(ctx)

	engine.Close()
}
`
