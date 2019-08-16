package web

import "microlog/web/actions"

func initRoutes() {
	router.GET("/", actions.Home)
	router.GET("/inputs", actions.Inputs)
	router.GET("/input/add", actions.AddInput)
	router.POST("/input/add", actions.CreateInput)
	router.POST("/input/stop/:id", actions.StopInput)
	router.POST("/input/start/:id", actions.StartInput)
	router.POST("/input/delete/:id", actions.DeleteInput)
	router.GET("/search", actions.Search)

	// static files
	router.Static("/assets", "./web/assets/public")
	router.StaticFile("/favicon.ico", "./web/assets/public/favicon.ico")
}
