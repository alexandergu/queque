package api

import (
	"net/http"

	"github.com/alexandergu/queque/internal/httpx"
)

func (router *Router) handleResizeWorkers(w http.ResponseWriter, r *http.Request) {
	data, err := httpx.Convert[ResizeWorkersData](r)
	if err != nil {
		httpx.Error(w, err)

		return
	}

	if err = router.engine.ResizeWorkersCount(data.Count); err != nil {
		httpx.Error(w, err)

		return
	}

	httpx.Ok(w)
}

func (router *Router) handleWorkersCount(w http.ResponseWriter, r *http.Request) {
	httpx.Resource(w, WorkersCount{
		Count: router.engine.WorkersCount(),
	}, func(s WorkersCount) any {
		return s
	})
}
