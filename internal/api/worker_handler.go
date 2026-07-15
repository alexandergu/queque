package api

import (
	"net/http"

	"github.com/alexandergu/queque/internal/httpx"
)

func (router *Router) handleResizeWorkers(w http.ResponseWriter, r *http.Request) {
	data, err := httpx.Convert[ResizeWorkersDto](r)
	if err != nil {
		httpx.Error(w, err)

		return
	}

	if err = router.engine.Resize(data.Count); err != nil {
		httpx.Error(w, err)
	}

	httpx.Ok(w)
}
