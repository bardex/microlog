package web

import "microlog/web/actions"

func initRoutes() {
	router.GET("/", actions.Home)
}
