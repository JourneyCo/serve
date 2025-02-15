package router

import "github.com/gorilla/mux"

type ServeRouter struct {
	*mux.Router
}

func (r ServeRouter) RBAC(int) ServeRouter {
	return r
}

func (r ServeRouter) SubPath(s string) ServeRouter {
	return ServeRouter{r.PathPrefix(s).Subrouter()}
}
