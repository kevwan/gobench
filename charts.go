package gobench

import (
	"net/http"
	"sort"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func generateLineItems(vals []Metrics) []opts.LineData {
	items := make([]opts.LineData, 0, len(vals))

	for _, v := range vals {
		items = append(items, opts.LineData{
			Value: int(v.P99 / time.Microsecond),
		})
	}

	return items
}

func generateChart(title string, bucket map[int]Metrics) http.HandlerFunc {
	keys := make([]int, 0, len(bucket))
	for k := range bucket {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	p50Vals := make([]opts.LineData, 0, len(bucket))
	for _, k := range keys {
		p50Vals = append(p50Vals, opts.LineData{Value: bucket[k].Median / time.Microsecond})
	}
	p99Vals := make([]opts.LineData, 0, len(bucket))
	for _, k := range keys {
		p99Vals = append(p99Vals, opts.LineData{Value: bucket[k].P99 / time.Microsecond})
	}

	return func(w http.ResponseWriter, _ *http.Request) {
		line := charts.NewLine()
		if len(title) > 0 {
			line.SetGlobalOptions(
				charts.WithTitleOpts(opts.Title{
					Title: title,
				}),
			)
		}
		line.SetGlobalOptions(
			charts.WithXAxisOpts(opts.XAxis{
				Name: "Time (s)",
			}),
			charts.WithYAxisOpts(opts.YAxis{
				Name: "Response Time (us)",
			}),
		)

		// Put data into instance
		line.SetXAxis(keys).AddSeries("P50", p50Vals).AddSeries("P99", p99Vals)
		line.Render(w)
	}
}
