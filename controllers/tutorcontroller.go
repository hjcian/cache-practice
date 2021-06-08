package controllers

import (
	"cachepractice/repository"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TutorController struct {
	DB     repository.Repository
	Logger *zap.Logger
}

func (t TutorController) List(c *gin.Context) {
	langSlug := c.Param("language_slug")
	tutors, err := t.DB.GetTutors(context.Background(), langSlug)
	if err != nil {
		if err == repository.ErrRecordNotFound {
			t.Logger.Warn("tutors not found", zap.String("langSlug", langSlug))
			c.JSON(http.StatusNotFound, gin.H{"error": "tutors not found"})
			return
		}
		t.Logger.Warn("list tutors error", zap.String("langSlug", langSlug))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "list tutors error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": tutors,
	})
}

func (t TutorController) Get(c *gin.Context) {
	tutorSlug := c.Param("tutor_slug")
	tutor, err := t.DB.GetTutor(context.Background(), tutorSlug)
	if err != nil {
		if err == repository.ErrRecordNotFound {
			t.Logger.Warn("tutor not found", zap.String("tutorSlug", tutorSlug))
			c.JSON(http.StatusNotFound, gin.H{"error": "tutors not found"})
			return
		}
		t.Logger.Warn("find tutor error", zap.String("tutorSlug", tutorSlug), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "find tutor error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": tutor,
	})
}
