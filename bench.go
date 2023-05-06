package gobench

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/zeromicro/go-zero/core/timex"
)

const defaultAddr = "localhost:8081"

type (
	Metrics struct {
		Median time.Duration
		P99    time.Duration
	}

	Bench struct {
		records   map[int]Metrics
		startTime time.Duration
		current   time.Duration
		bucket    taskHeap
	}
)

func NewBench() *Bench {
	return &Bench{
		records:   make(map[int]Metrics),
		startTime: timex.Now(),
		current:   timex.Now(),
	}
}

func (b *Bench) Run(qps int, fn func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	ticket := time.NewTicker(time.Second / time.Duration(qps))
	defer ticket.Stop()

	for {
		select {
		case <-ticket.C:
			b.runSingle(fn)
		case <-c:
			signal.Stop(c)
			go func() {
				time.Sleep(time.Second)
				openBrowser("http://" + defaultAddr)
			}()
			goto chart
		}
	}

chart:
	http.HandleFunc("/", generateChart(b.records))
	http.ListenAndServe(defaultAddr, nil)
}

func (b *Bench) runSingle(fn func()) {
	start := timex.Now()
	fn()
	elapsed := timex.Since(start)

	if timex.Since(b.current) > time.Second {
		metrics := calculate(b.bucket)
		index := int((b.current - b.startTime) / time.Second)
		b.records[index] = metrics
		b.current = start
		b.bucket = taskHeap{}
	}

	b.bucket.Push(Task{
		Duration: elapsed,
	})
}

func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		fmt.Println(err)
	}
}
