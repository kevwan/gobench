package gobench

import (
	"container/heap"
	"time"
)

type (
	Task struct {
		time.Duration
	}

	taskHeap []Task
)

func (h *taskHeap) Len() int {
	return len(*h)
}

func (h *taskHeap) Less(i, j int) bool {
	return (*h)[i].Duration < (*h)[j].Duration
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
		} else if (*h)[0].Duration < each.Duration {
			heap.Pop(h)
			heap.Push(h, each)
		}
	}

	return *h
}

func calculate(bucket taskHeap) Metrics {
	var metrics Metrics
	size := bucket.Len()
	if size == 0 {
		return metrics
	}

	fiftyPercent := size >> 1
	if fiftyPercent > 0 {
		top50pTasks := topK(bucket, fiftyPercent)
		medianTask := top50pTasks[0]
		metrics.Median = medianTask.Duration
		onePercent := fiftyPercent / 50
		if onePercent > 0 {
			top1pTasks := topK(top50pTasks, onePercent)
			task99th := top1pTasks[0]
			metrics.P99 = task99th.Duration
		} else {
			mostDuration := getTopDuration(top50pTasks)
			metrics.P99 = mostDuration
		}
	} else {
		mostDuration := getTopDuration(bucket)
		metrics.Median = mostDuration
		metrics.P99 = mostDuration
	}

	return metrics
}

func getTopDuration(tasks []Task) time.Duration {
	top := topK(tasks, 1)
	if len(top) < 1 {
		return 0
	}

	return top[0].Duration
}
