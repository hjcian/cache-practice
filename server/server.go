package server

import (
	"cachepractice/controllers"
	"cachepractice/repository"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func NewRouter(db repository.Repository, logger *zap.Logger) *gin.Engine {
	router := gin.Default()
	router.HandleMethodNotAllowed = true

	tutor := controllers.TutorController{
		DB:     db,
		Logger: logger,
	}
	router.GET("/api/tutors/:language_slug", tutor.List)
	router.GET("/api/tutor/:tutor_slug", tutor.Get)

	return router
}
