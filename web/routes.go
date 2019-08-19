package web

import "microlog/web/actions"

func initRoutes() {
	router.GET("/", actions.Home)
	router.GET("/api/inputs", actions.Inputs)
	router.GET("/input/add", actions.AddInput)
	router.POST("/input/add", actions.CreateInput)
	router.POST("/api/input/stop/:id", actions.StopInput)
	router.POST("/api/input/start/:id", actions.StartInput)
	router.POST("/api/input/delete/:id", actions.DeleteInput)
	router.GET("/api/search", actions.Search)

	// static files
	router.Static("/assets", "./web/assets/public")
	router.StaticFile("/favicon.ico", "./web/assets/public/favicon.ico")
}
