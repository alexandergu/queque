package api

import (
	"net/http"

	"github.com/alexandergu/queque/internal/queue"
)

type Router struct {
	mux    *http.ServeMux
	engine *queue.Engine
}

func NewRouter(engine *queue.Engine) *Router {
	router := &Router{
		mux:    http.NewServeMux(),
		engine: engine,
	}
	router.initRoutes()

	return router
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.mux.ServeHTTP(w, r)
}

func (router *Router) initRoutes() {
	router.mux.HandleFunc("POST /api/jobs", router.handleCreateJob)
	router.mux.HandleFunc("GET /api/jobs", router.handleGetAllJobs)
	router.mux.HandleFunc("GET /api/jobs/{job}", router.handleGetJob)
}
