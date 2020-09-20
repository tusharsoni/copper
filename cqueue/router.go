package cqueue

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tusharsoni/copper/chttp"
	"github.com/tusharsoni/copper/clogger"
)

type Router struct {
	rw     chttp.ReaderWriter
	logger clogger.Logger

	svc Svc
}

type RouterParams struct {
	RW     chttp.ReaderWriter
	Logger clogger.Logger

	Svc Svc
}

func NewRouter(p RouterParams) chttp.Router {
	ro := &Router{
		rw:     p.RW,
		logger: p.Logger,
		svc:    p.Svc,
	}

	return chttp.NewRouter([]chttp.Route{
		{
			Path:    "/api/queue/tasks/{uuid}",
			Methods: []string{http.MethodGet},
			Handler: http.HandlerFunc(ro.HandleGetTask),
		},
	})
}

func (ro *Router) HandleGetTask(w http.ResponseWriter, r *http.Request) {
	taskUUID := mux.Vars(r)["uuid"]

	task, err := ro.svc.GetTask(r.Context(), taskUUID)
	if err != nil {
		ro.logger.Error("Failed to get task", err)
		ro.rw.InternalErr(w)
		return
	}

	ro.rw.OK(w, task)
}
