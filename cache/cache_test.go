package cache

import (
	"cachepractice/repository"
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	gocache "github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

const waitToJobDone = time.Second

type recorder struct {
	GetTutorCount  int
	GetTutorsCount int
	mutex          sync.Mutex
}

func (m *recorder) GetTutor(ctx context.Context, tutorSlug string) (*repository.TutorInfo, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.GetTutorCount++
	time.Sleep(waitToJobDone)
	return &englishTutor1, nil
}

func (m *recorder) GetTutors(ctx context.Context, langSlug string) ([]repository.TutorInfo, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.GetTutorsCount++
	time.Sleep(waitToJobDone)
	return []repository.TutorInfo{englishTutor1}, nil
}

func Test_cacheLayer_GetTutors_Cache_Miss_Should_Block(t *testing.T) {
	// t.Skip()

	logger, _ := zap.NewDevelopment()
	dbRecorder := recorder{}

	cache := cacheLayer{
		db:     &dbRecorder,
		logger: logger,
		cache:  gocache.New(gocache.NoExpiration, gocache.NoExpiration),
		ml:     NewMultipleLock(),
	}

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cache.GetTutors(context.Background(), "foo")
		}()
	}
	wg.Wait()
	assert.Equal(t, dbRecorder.GetTutorsCount, 1)
}

func Test_cacheLayer_GetTutors_Cache_Expired_Should_Update_Once(t *testing.T) {
	// t.Skip()

	logger, _ := zap.NewDevelopment()
	dbRecorder := recorder{}
	lang := "english"
	cacheKey := fmt.Sprintf(langCacheKey, lang)

	cache := cacheLayer{
		db:     &dbRecorder,
		logger: logger,
		cache:  gocache.New(gocache.NoExpiration, gocache.NoExpiration),
		ml:     NewMultipleLock(),
	}
	// create an expired cache data
	expiredTime := time.Now().Add(-24 * time.Hour)
	cache.updateCache(cacheKey, []repository.TutorInfo{oldEnglishTutor1}, expiredTime)
	assert.Equal(t, 0, dbRecorder.GetTutorsCount)

	// first call to get old data and trigger the async job
	data, err := cache.GetTutors(context.Background(), lang)
	assert.NoError(t, err)
	assert.Equal(t, []repository.TutorInfo{oldEnglishTutor1}, data)

	// wait for the async job done
	time.Sleep(2 * waitToJobDone)
	data, err = cache.GetTutors(context.Background(), lang)
	assert.NoError(t, err)
	assert.Equal(t, []repository.TutorInfo{englishTutor1}, data)
	assert.Equal(t, 1, dbRecorder.GetTutorsCount)
}

func Test_cacheLayer_GetTutor_Cache_Miss_Should_Block(t *testing.T) {
	// t.Skip()

	logger, _ := zap.NewDevelopment()
	dbRecorder := recorder{}

	cache := cacheLayer{
		db:     &dbRecorder,
		logger: logger,
		cache:  gocache.New(gocache.NoExpiration, gocache.NoExpiration),
		ml:     NewMultipleLock(),
	}

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cache.GetTutor(context.Background(), "foo")
		}()
	}
	wg.Wait()
	assert.Equal(t, dbRecorder.GetTutorCount, 1)
}

func Test_cacheLayer_GetTutor_Cache_Expired_Should_Update_Once(t *testing.T) {
	// t.Skip()

	logger, _ := zap.NewDevelopment()
	dbRecorder := recorder{}
	tutorSlug := oldEnglishTutor1.TutorSlug
	cacheKey := fmt.Sprintf(tutorCacheKey, tutorSlug)

	cache := cacheLayer{
		db:     &dbRecorder,
		logger: logger,
		cache:  gocache.New(gocache.NoExpiration, gocache.NoExpiration),
		ml:     NewMultipleLock(),
	}
	// create an expired cache data
	expiredTime := time.Now().Add(-24 * time.Hour)
	cache.updateCache(cacheKey, oldEnglishTutor1, expiredTime)
	assert.Equal(t, 0, dbRecorder.GetTutorCount)

	// first call to get old data and trigger the async job
	data, err := cache.GetTutor(context.Background(), tutorSlug)
	assert.NoError(t, err)
	assert.Equal(t, oldEnglishTutor1, *data)

	// wait for the async job done
	time.Sleep(2 * waitToJobDone)
	data, err = cache.GetTutor(context.Background(), tutorSlug)
	assert.NoError(t, err)
	assert.Equal(t, englishTutor1, *data)
	assert.Equal(t, dbRecorder.GetTutorCount, 1)
}

var (
	englishTutor1 = repository.TutorInfo{
		TutorID:       "1",
		TutorSlug:     "foo-bar",
		TutorName:     "Amazing Teacher 1",
		TutorHeadline: "Hi I'm a English Teacher",
		PriceInfo: repository.PriceInfo{
			Trial:  5,
			Normal: 10,
		},
		TeachingLangs: []int{123, 321},
	}
	oldEnglishTutor1 = repository.TutorInfo{
		TutorID:       "1",
		TutorSlug:     "foo-bar",
		TutorName:     "Amazing Teacher 2",
		TutorHeadline: "Hi I'm a English Teacher",
		PriceInfo: repository.PriceInfo{
			Trial:  100,
			Normal: 200,
		},
		TeachingLangs: []int{123, 321},
	}
)
