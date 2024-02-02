package router

import "app/internal/api/controllers"

func (router *Router) ShorterRouter() {
	router.engine.POST("/shorter", controllers.DoShorter)
	router.engine.GET("short/:key", controllers.ShortRedirect)
}
