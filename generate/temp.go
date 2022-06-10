package generate

import (
	"fmt"
	"io"
	"os"
	"regexp"
)

const (
	version = "v1.0.1-alpha"
)

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

func SetProjectModByFile(modPath string) error {
	f, err := os.Open(modPath)
	if err != nil {
		return err
	}

	content, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	mods := regexp.MustCompile(`module (.+)`).FindAllStringSubmatch(string(content), -1)

	if len(mods) == 0 {
		return fmt.Errorf("this is not a valid go mod file, path: %s", modPath)
	}

	SetProjectMod(mods[0][1])
	return nil
}

var spiderTemp = `package spiders

import (
	"fmt"
	"sync"

	"github.com/wheat-os/slubby/spider"
	"github.com/wheat-os/slubby/stream"
)

type {{.name}}Spider struct{}

func (t *{{.name}}Spider) UId() string {
	return "{{.uid}}"
}

func (t *{{.name}}Spider) FQDN() string {
	return "{{.fqdn}}"
}

func (t *{{.name}}Spider) Parse(response *stream.HttpResponse) (stream.Stream, error) {
	fmt.Println(response.Text())

	return nil, nil
}

func (t *{{.name}}Spider) StartRequest() stream.Stream {
	req, _ := stream.Request(t, "http://www.baidu.com", nil)
	return req
}

var (
	{{.name}}Once = sync.Once{}
	{{.name}}P *{{.name}}Spider
)

func {{.nameCap}}Spider() spider.Spider {
	{{.name}}Once.Do(func() {
		{{.name}}P = &{{.name}}Spider{}
	})
	return {{.name}}P
}
`

var settingTemp = `package {{.pkgName}}

import (
	"os"

	"github.com/wheat-os/slubby/download"
	"github.com/wheat-os/slubby/download/middle"
	"github.com/wheat-os/slubby/engine"
	"github.com/wheat-os/slubby/outputter"
	"github.com/wheat-os/slubby/outputter/pipline"
	"github.com/wheat-os/slubby/scheduler"
	"github.com/wheat-os/slubby/scheduler/buffer"
	"github.com/wheat-os/slubby/scheduler/filter"
	log "github.com/wheat-os/wlog"
)

// ***************************************** Logger *****************************************
func init() {
	log.SetStdOptions(log.WithDisPlayLevel(log.InfoLevel))
	log.SetStdOptions(log.WithDisableCaller(true))
	log.SetStdOptions(log.WithOutput(os.Stdout))

	// set the log level
	log.SetStdOptions(log.WithDisPlayLevel(log.DebugLevel))
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

// http headers
var defaultHttpRequestHeaders = map[string]string{
	"user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
}

var downloadMiddlewareModule = middle.MiddleGroup(
	middle.HeadMiddle(defaultHttpRequestHeaders),
	middle.LogMiddle(),
)

// download limiter
// var downloadLimiterModule = limiter.ShortLimiter(
// 	time.Second * 3,
// )

var downloadModule = download.ShortDownload(
	// download.WithLimiter(downloadLimiterModule),
	download.WithDownloadMiddle(downloadMiddlewareModule),

	// the number of downloader retries
	download.WithDownloadRetry(2),
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
version = "{{.version}}"

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

var registerTemp = `package main

import (
	"{{.projectMod}}/{{.projectName}}"

	"{{.projectMod}}/spiders"
)

func init() {
	{{- range $index, $value := .spiders}}
		{{$.projectName}}.DefaultEngine.Register(spiders.{{$value}}())
	{{- end}}
}
`
