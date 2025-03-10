package routes

import (
	"library_project/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/songs", handlers.GetSongText)
	r.GET("/songs/:id/text", handlers.GetSongText)
	r.POST("/songs", handlers.AddSong)
	r.PUT("/songs/:id", handlers.UpdateSong)
	r.DELETE("/songs/:id", handlers.DeleteSong)
}
