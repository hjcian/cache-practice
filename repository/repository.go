package repository

import (
	"context"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Repository interface {
	GetTutors(ctx context.Context, langSlug string) ([]TutorInfo, error)
	GetTutor(ctx context.Context, tutorSlug string) (*TutorInfo, error)
}

type PriceInfo struct {
	Trial  float32 `json:"trial"`
	Normal float32 `json:"normal"`
}

type TutorInfo struct {
	TutorID       string    `json:"id"`
	TutorSlug     string    `json:"slug"`
	TutorName     string    `json:"name"`
	TutorHeadline string    `json:"headline"`
	PriceInfo     PriceInfo `json:"price_info"`
	TeachingLangs []int     `json:"teaching_languages"`
}
