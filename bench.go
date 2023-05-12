package gobench

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/zeromicro/go-zero/core/cmdline"
	"github.com/zeromicro/go-zero/core/threading"
	"github.com/zeromicro/go-zero/core/timex"
)

const (
	defaultHost             = "localhost"
	defaultPath             = "/"
	minCpusToReserveForCalc = 4
)

type (
	Bench struct {
		records   map[int]metrics
		startTime time.Duration
		current   time.Duration
		bucket    []Task
		title     string
		host      string
		port      int
		duration  time.Duration
		ticker    *time.Ticker
		lock      sync.Mutex
		quit      chan struct{}
	}

	Config struct {
		Host     string
		Port     int
		Title    string
		Duration time.Duration
	}

	metrics struct {
		Average   time.Duration
		P50       time.Duration
		P90       time.Duration
		P99       time.Duration
		Qps       int
		ErrorRate float64
		Cpu       float64
		Memory    float64
	}

	chartConfig struct {
		title         string
		plotErrorRate bool
	}
)

func NewBench() *Bench {
	return &Bench{
		records:   make(map[int]metrics),
		startTime: timex.Now(),
		current:   timex.Now(),
		quit:      make(chan struct{}),
	}
}

func NewBenchWithConfig(cfg Config) *Bench {
	return &Bench{
		records:   make(map[int]metrics),
		startTime: timex.Now(),
		current:   timex.Now(),
		title:     cfg.Title,
		host:      cfg.Host,
		port:      cfg.Port,
		duration:  cfg.Duration,
		quit:      make(chan struct{}),
	}
}

func (b *Bench) Run(qps int, fn func()) {
	b.doRun(qps, func() error {
		fn()
		return nil
	}, false)
}

func (b *Bench) RunErr(qps int, fn func() error) {
	b.doRun(qps, fn, true)
}

func (b *Bench) doRun(qps int, fn func() error, hasErr bool) {
	fmt.Println("Ctrl+C to show the benchmark charts")

	interval := time.Second / time.Duration(qps)
	b.ticker = time.NewTicker(interval)
	b.runLoop(fn)
	b.ticker.Stop()

	listener, err := net.Listen("tcp", b.buildAddr())
	if err != nil {
		fmt.Println(err)
		return
	}

	addr := listener.Addr().String()
	fmt.Printf("\n\nListening on: %s\n", addr)

	go func() {
		http.HandleFunc(defaultPath, generateChart(b.records, chartConfig{
			title:         b.title,
			plotErrorRate: hasErr,
		}))
		if err := http.Serve(listener, nil); err != nil {
			fmt.Println(err)
		}
	}()

	time.Sleep(time.Millisecond * 500)
	openBrowser("http://" + addr)
	cmdline.EnterToContinue()
}

func (b *Bench) analyze() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	var seconds int
	// discard last second before stop, so we can get a more accurate result
	// because the last second may not be a complete second
	for {
		select {
		case <-b.quit:
			return
		case <-ticker.C:
			if seconds%60 == 0 {
				fmt.Printf("\n%2dm ", seconds/60)
			}
			fmt.Print(".")
			seconds++

			var bucket []Task

			b.lock.Lock()
			index := int((b.current - b.startTime) / time.Second)
			b.current = timex.Now()
			for i, task := range b.bucket {
				if timex.Since(task.start) < time.Second {
					bucket = b.bucket[:i]
					b.bucket = b.bucket[i:]
					break
				}
			}
			b.lock.Unlock()

			if len(bucket) == 0 {
				continue
			}

			metrics := calculate(bucket)
			metrics.Cpu = getCpuUsage()
			metrics.Memory = getMemoryUsage()

			b.lock.Lock()
			b.records[index] = metrics
			b.lock.Unlock()
		}
	}
}

func (b *Bench) buildAddr() string {
	host := b.host
	if len(host) == 0 {
		host = defaultHost
	}

	return fmt.Sprintf("%s:%d", host, b.port)
}

func (b *Bench) collect(collector <-chan Task) {
	for {
		select {
		case <-b.quit:
			return
		case task := <-collector:
			b.lock.Lock()
			b.bucket = append(b.bucket, task)
			b.lock.Unlock()
		}
	}
}

func (b *Bench) runLoop(fn func() error) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	channel := make(chan struct{}, 1)
	collector := make(chan Task, 1)
	numCpus := runtime.NumCPU()
	if numCpus >= minCpusToReserveForCalc {
		// keep one cpu for metrics calculation
		numCpus -= 1
	}
	group := threading.NewWorkerGroup(func() {
		for {
			select {
			case <-b.quit:
				return
			case <-channel:
				b.runSingle(collector, fn)
			}
		}
	}, numCpus)
	go group.Start()
	go b.collect(collector)
	go b.analyze()

	stop := func() {
		signal.Stop(c)
		close(b.quit)
	}

	timer := time.NewTimer(b.duration)
	defer timer.Stop()

	for {
		select {
		case <-b.ticker.C:
			channel <- struct{}{}
		case <-timer.C:
			stop()
			return
		case <-c:
			stop()
			return
		}
	}
}

func (b *Bench) runSingle(collector chan<- Task, fn func() error) {
	start := timex.Now()
	err := fn()
	duration := timex.Since(start)
	collector <- Task{
		start:    start,
		duration: duration,
		err:      err,
	}
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
