package queue

type JobHeap []*Job

func (h JobHeap) Less(i, j int) bool {
	if h[i].Priority == h[j].Priority {
		return h[i].seq < h[j].seq
	}

	return h[i].Priority > h[j].Priority
}

func (h JobHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *JobHeap) Push(x any) {
	*h = append(*h, x.(*Job))
}

func (h *JobHeap) Pop() any {
	heap := *h
	length := len(heap)
	item := heap[length-1]

	heap[length-1] = nil
	*h = heap[:length-1]

	return item
}

func (h JobHeap) Len() int {
	return len(h)
}
