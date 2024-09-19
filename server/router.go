package server

import (
	"fmt"
	"hireme-api/config"
	"hireme-api/controllers"
	"hireme-api/middleware"
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

	if config.GetEnv("ENVIRONMENT") == "development" {
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
	teapot := new(controllers.TeapotController)
	hint := new(controllers.HintController)
	score := new(controllers.ScoreController)

	api := router.Group("api")
	{
		v69 := api.Group("v69")
		{
			v69.GET("/teapot", teapot.Teapot)
			v69.POST("/teapot", teapot.Teapot)
			v69.POST("/bringiton", public.Init)
			v69.Use(middleware.LoadSession())
			{
				v69.POST("/whatsupdoc", public.Status)
				v69.POST("/answer", public.Answer)
				v69.POST("/hints", hint.GetHints)
				v69.POST("/takehint", hint.Hint)
				v69.POST("/scores", score.Scores)
			}
		}
	}
	return router

}
