package server

import (
	"log"

	"github.com/4Noyis/url-shortener/internal/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(urlHandler *handlers.URLHandler) *gin.Engine {
	router := gin.Default()

	api := router.Group("/api/v1")
	{
		data := api.Group("/data")
		{
			data.POST("/shorten", urlHandler.ShortenURL)
		}
		
		api.GET("/:shortURL", urlHandler.RedirectURL)
	}

	return router
}

func StartServer(router *gin.Engine, port string) {
	log.Printf("Starting server on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}