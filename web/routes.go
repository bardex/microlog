package web

import "microlog/web/actions"

func initRoutes() {
	router.GET("/", actions.Home)

	// static files
	router.Static("/assets", "./web/assets")
	router.StaticFile("/favicon.ico", "./web/assets/favicon.ico")
}
