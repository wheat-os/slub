package douban

import (
	"os"
	"time"

	"github.com/wheat-os/slubby/download"
	"github.com/wheat-os/slubby/download/limiter"
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
var downloadLimiterModule = limiter.ShortLimiter(
	time.Millisecond * 300,
)

var downloadModule = download.ShortDownload(
	download.WithLimiter(downloadLimiterModule),
	download.WithDownloadMiddle(downloadMiddlewareModule),

	// the number of downloader retries
	download.WithDownloadRetry(2),
)

// *************************************** Outputter **************************************

// pipline
var piplineModule = pipline.GroupPipline(
	pipline.SaveCsvPipline("./data.csv", 0),
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
