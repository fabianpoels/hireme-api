package server

import (
	"fmt"
	"hireme-api/config"
	"hireme-api/controllers"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// FOR DEV PURPOSES!!!!
	// TODO: rework to dynamically configure cors, depending on the env
	corsConfig := cors.DefaultConfig()

	if config.GetEnv("ENVIRONMENT") == "dev" {
		// LOCAL DEV CONFIG
		// domain := config.GetEnv("DOMAIN")
		router.SetTrustedProxies(nil)
		corsConfig.AllowOrigins = []string{"http://localhost:5173", "http://127.0.0.1:5173"}
		corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control"}
		corsConfig.AllowMethods = []string{"POST", "GET", "PUT", "OPTIONS"}
		corsConfig.ExposeHeaders = []string{"Content-Length"}
		corsConfig.AllowCredentials = true
		corsConfig.MaxAge = 12 * time.Hour
	} else {
		domain := config.GetEnv("DOMAIN")
		router.SetTrustedProxies(nil)
		corsConfig.AllowOrigins = []string{fmt.Sprintf("http://%s", domain)}
		corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control"}
		corsConfig.AllowMethods = []string{"POST", "GET", "PUT", "OPTIONS"}
		corsConfig.ExposeHeaders = []string{"Content-Length"}
		corsConfig.AllowCredentials = true
		corsConfig.MaxAge = 12 * time.Hour
	}

	router.Use(cors.New(corsConfig))

	public := new(controllers.PublicController)

	api := router.Group("api")
	{
		v69 := api.Group("v69")
		{
			// public routes
			v69.POST("/whatsupdoc", public.Status)
			v69.POST("/bringiton", public.Init)
		}
	}
	return router

}
