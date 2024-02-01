package router

import "app/internal/api/controllers"

func (router *Router) MainRoutes() {
	router.engine.GET("/", controllers.GetMain)
}
