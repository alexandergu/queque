package api

import (
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
