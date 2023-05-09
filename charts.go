package gobench

import (
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func generateChart(title string, bucket map[int]Metrics) http.HandlerFunc {
	keys := make([]int, 0, len(bucket))
	for k := range bucket {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	avgVals := make([]opts.LineData, 0, len(bucket))
	for _, k := range keys {
		avgVals = append(avgVals, opts.LineData{Value: bucket[k].Average / time.Microsecond})
	}
	p50Vals := make([]opts.LineData, 0, len(bucket))
	for _, k := range keys {
		p50Vals = append(p50Vals, opts.LineData{Value: bucket[k].P50 / time.Microsecond})
	}
	p90Vals := make([]opts.LineData, 0, len(bucket))
	for _, k := range keys {
		p90Vals = append(p90Vals, opts.LineData{Value: bucket[k].P90 / time.Microsecond})
	}
	p99Vals := make([]opts.LineData, 0, len(bucket))
	for _, k := range keys {
		p99Vals = append(p99Vals, opts.LineData{Value: bucket[k].P99 / time.Microsecond})
	}

	respTimeLine := charts.NewLine()
	if len(title) > 0 {
		respTimeLine.SetGlobalOptions(
			charts.WithTitleOpts(opts.Title{
				Title: title,
			}),
		)
	}
	respTimeLine.SetGlobalOptions(
		charts.WithXAxisOpts(opts.XAxis{
			Name: "Time (s)",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: "Response Time (us)",
		}),
	)
	respTimeLine.
		SetXAxis(keys).
		AddSeries("Average", avgVals).
		AddSeries("P50", p50Vals).
		AddSeries("P90", p90Vals).
		AddSeries("P99", p99Vals)

	qpsLine := charts.NewLine()
	qpsLine.SetGlobalOptions(
		charts.WithXAxisOpts(opts.XAxis{
			Name: "Time (s)",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: "QPS",
		}),
	)

	qpsVals := make([]opts.LineData, 0, len(bucket))
	for _, k := range keys {
		qpsVals = append(qpsVals, opts.LineData{Value: bucket[k].Qps})
	}
	qpsLine.SetXAxis(keys).AddSeries("QPS", qpsVals)

	cpuLine := charts.NewLine()
	cpuLine.SetGlobalOptions(
		charts.WithXAxisOpts(opts.XAxis{
			Name: "Time (s)",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: "Percent(%)",
		}),
	)

	cpuVals := make([]opts.LineData, 0, len(bucket))
	for _, k := range keys {
		cpuVals = append(cpuVals, opts.LineData{Value: bucket[k].Cpu})
	}
	cpuLine.SetXAxis(keys).AddSeries("CPU", cpuVals)

	memLine := charts.NewLine()
	memLine.SetGlobalOptions(
		charts.WithXAxisOpts(opts.XAxis{
			Name: "Time (s)",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: "Percent(%)",
		}),
	)

	memVals := make([]opts.LineData, 0, len(bucket))
	for _, k := range keys {
		memVals = append(memVals, opts.LineData{Value: bucket[k].Memory})
	}
	memLine.SetXAxis(keys).AddSeries("Memory", memVals)

	return func(w http.ResponseWriter, _ *http.Request) {
		if err := respTimeLine.Render(w); err != nil {
			fmt.Println(err)
		}

		if err := qpsLine.Render(w); err != nil {
			fmt.Println(err)
		}

		if err := cpuLine.Render(w); err != nil {
			fmt.Println(err)
		}

		if err := memLine.Render(w); err != nil {
			fmt.Println(err)
		}
	}
}
