package web

import "microlog/web/actions"

func initRoutes() {
	router.GET("/", actions.Home)
	router.GET("/inputs", actions.Inputs)
	router.GET("/search", actions.Search)

	// static files
	router.Static("/assets", "./web/assets")
	router.StaticFile("/favicon.ico", "./web/assets/favicon.ico")
}
