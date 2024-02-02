package router

import (
	_ "app/docs"
	"app/internal/api/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"os"
	"path/filepath"
	"runtime"
)

type Router struct {
	engine *gin.Engine
}

func NewRouter() Router {
	return Router{engine: gin.Default()}
}

func (router *Router) Run(port string) error {
	router.Setup()
	return router.engine.Run(":" + port)
}

func (router *Router) Setup() {
	gin.SetMode(gin.DebugMode)
	var domain string
	if os.Getenv("SETUP_TYPE") == "local" {
		domain = "localhost"
	} else {
		domain = "shorter.ultraevs.ru"
	}
	router.engine.Use(middleware.AuthMiddleware(domain))
	_, currentFilePath, _, _ := runtime.Caller(1)
	templatesPath := filepath.Join(filepath.Dir(currentFilePath), "../../../templates")
	router.engine.LoadHTMLGlob(filepath.Join(templatesPath, "*.html"))
	router.engine.Static("/assets", filepath.Join(templatesPath, "assets"))
	router.engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))
	router.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.MainRoutes()
	router.ShorterRouter()
}
