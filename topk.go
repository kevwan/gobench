package gobench

import (
	"container/heap"
	"time"
)

type (
	Task struct {
		start    time.Duration
		duration time.Duration
	}

	taskHeap []Task
)

func (h *taskHeap) Len() int {
	return len(*h)
}

func (h *taskHeap) Less(i, j int) bool {
	return (*h)[i].duration < (*h)[j].duration
}

func (h *taskHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *taskHeap) Push(x any) {
	*h = append(*h, x.(Task))
}

func (h *taskHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func topK(all []Task, k int) []Task {
	h := new(taskHeap)
	heap.Init(h)

	for _, each := range all {
		if h.Len() < k {
			heap.Push(h, each)
		} else if (*h)[0].duration < each.duration {
			heap.Pop(h)
			heap.Push(h, each)
		}
	}

	return *h
}

func calculate(bucket []Task) Metrics {
	var metrics Metrics
	size := len(bucket)
	if size == 0 {
		return metrics
	}

	var total time.Duration
	for _, each := range bucket {
		total += each.duration
	}
	metrics.Average = total / time.Duration(size)

	fiftyPercent := size >> 1
	if fiftyPercent > 0 {
		top50pTasks := topK(bucket, fiftyPercent)
		medianTask := top50pTasks[0]
		metrics.P50 = medianTask.duration
		tenPercent := fiftyPercent / 5
		if tenPercent > 0 {
			top10pTasks := topK(top50pTasks, tenPercent)
			task90th := top10pTasks[0]
			metrics.P90 = task90th.duration
			onePercent := tenPercent / 10
			if onePercent > 0 {
				top1pTasks := topK(top10pTasks, onePercent)
				task99th := top1pTasks[0]
				metrics.P99 = task99th.duration
			} else {
				mostDuration := getTopDuration(top50pTasks)
				metrics.P99 = mostDuration
			}
		} else {
			mostDuration := getTopDuration(top50pTasks)
			metrics.P90 = mostDuration
			metrics.P99 = mostDuration
		}
	} else {
		mostDuration := getTopDuration(bucket)
		metrics.P50 = mostDuration
		metrics.P90 = mostDuration
		metrics.P99 = mostDuration
	}

	metrics.Qps = size

	return metrics
}

func getTopDuration(tasks []Task) time.Duration {
	top := topK(tasks, 1)
	if len(top) < 1 {
		return 0
	}

	return top[0].duration
}
