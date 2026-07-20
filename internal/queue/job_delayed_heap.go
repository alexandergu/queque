package queue

type JobDelayedHeap []*Job

func (j JobDelayedHeap) Len() int {
	return len(j)
}

func (j JobDelayedHeap) Less(i, k int) bool {
	if j[i].AvailableAt.Equal(j[k].AvailableAt) {
		return j[i].seq < j[k].seq
	}

	return j[i].AvailableAt.Before(j[k].AvailableAt)
}

func (j JobDelayedHeap) Swap(i, k int) {
	j[i], j[k] = j[k], j[i]
}

func (j *JobDelayedHeap) Push(x any) {
	*j = append(*j, x.(*Job))
}

func (j *JobDelayedHeap) Pop() any {
	heap := *j
	length := len(heap)
	item := heap[length-1]

	heap[length-1] = nil
	*j = heap[:length-1]

	return item
}
