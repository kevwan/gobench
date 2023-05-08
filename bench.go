package gobench

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/zeromicro/go-zero/core/cmdline"
	"github.com/zeromicro/go-zero/core/timex"
)

const (
	defaultHost = "localhost"
	defaultPath = "/"
)

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
		title     string
		host      string
		port      int
	}

	Config struct {
		Host  string
		Port  int
		Title string
	}
)

func NewBench() *Bench {
	return &Bench{
		records:   make(map[int]Metrics),
		startTime: timex.Now(),
		current:   timex.Now(),
	}
}

func NewBenchWithConfig(cfg Config) *Bench {
	return &Bench{
		records:   make(map[int]Metrics),
		startTime: timex.Now(),
		current:   timex.Now(),
		title:     cfg.Title,
		host:      cfg.Host,
		port:      cfg.Port,
	}
}

func (b *Bench) Run(qps int, fn func()) {
	fmt.Println("Ctrl+C to show the benchmark chart")

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	ticket := time.NewTicker(time.Second / time.Duration(qps))
	defer ticket.Stop()

	for {
		select {
		case <-ticket.C:
			b.runSingle(fn)
		case <-c:
			signal.Stop(c)
			goto chart
		}
	}

chart:
	listener, err := net.Listen("tcp", b.buildAddr())
	if err != nil {
		fmt.Println(err)
		return
	}

	addr := listener.Addr().String()
	fmt.Printf("\nListening on: %s\nPress Enter to quit\n", addr)

	go func() {
		http.HandleFunc(defaultPath, generateChart(b.title, b.records))
		if err := http.Serve(listener, nil); err != nil {
			fmt.Println(err)
		}
	}()

	time.Sleep(time.Millisecond * 500)
	openBrowser("http://" + addr)
	cmdline.EnterToContinue()
}

func (b *Bench) buildAddr() string {
	host := b.host
	if len(host) == 0 {
		host = defaultHost
	}

	return fmt.Sprintf("%s:%d", host, b.port)
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
