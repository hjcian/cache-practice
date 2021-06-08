package cache

import (
	"cachepractice/repository"
	"context"
)

type CacheLayer struct {
	db repository.Repository
}

func (c CacheLayer) GetTutors(ctx context.Context, lang string) ([]repository.TutorInfo, error) {
	return c.db.GetTutors(ctx, lang)
}
func (c CacheLayer) GetTutor(ctx context.Context, tutorSlug string) (*repository.TutorInfo, error) {
	return c.db.GetTutor(ctx, tutorSlug)
}
