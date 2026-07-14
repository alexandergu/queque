package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/alexandergu/queque/internal/httpx"
)

func (router *Router) handleCreateJob(w http.ResponseWriter, r *http.Request) {
	data, err := httpx.Convert[CreateJobData](r)

	if err != nil {
		httpx.Error(w, err)

		return
	}

	job, _ := router.engine.Run(data.ToJobDto())
	httpx.Resource(w, job, RenderJob)
}

func (router *Router) handleGetAllJobs(w http.ResponseWriter, r *http.Request) {
	jobs := router.engine.GetJobs()

	httpx.Resources(w, jobs, RenderJob)
}

func (router *Router) handleGetJob(w http.ResponseWriter, r *http.Request) {
	job, err := GetJobFromRequest(router.engine, r)
	if err != nil {
		httpx.Error(w, err)

		return
	}

	httpx.Resource(w, job, RenderJob)
}

func (router *Router) handleCancelJob(w http.ResponseWriter, r *http.Request) {
	job, err := GetJobFromRequest(router.engine, r)
	if err != nil {
		httpx.Error(w, err)

		return
	}

	snapshot, err := router.engine.Cancel(job.ID)
	if err != nil {
		httpx.Error(w, err)

		return
	}

	httpx.Resource(w, snapshot, RenderJob)
}

func (router *Router) handleJobEvents(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)

	if !ok {
		httpx.Error(w, fmt.Errorf("streaming unsupported"))

		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")

	ch, unsubscribe := router.engine.Subscribe()
	defer unsubscribe()

	for {
		select {
		case <-r.Context().Done():
			return
		case event, ok := <-ch:
			if !ok {
				return
			}

			data, err := json.Marshal(event.Job)
			if err != nil {
				return
			}

			fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event.Type, data)

			flusher.Flush()
		}
	}
}
