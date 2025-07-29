package app

import "github.com/mingtmt/book-store/internal/routes"

type BaseModule struct {
	routes routes.Route
}

func (bm *BaseModule) Routes() routes.Route {
	return bm.routes
}
