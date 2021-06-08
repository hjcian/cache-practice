package cache

import (
	"errors"
	"fmt"
	"time"

	gocache "github.com/patrickmn/go-cache"
	"go.uber.org/zap"

	"cachepractice/repository"
	"context"
)

var (
	ErrCacheMissed     = errors.New("cache missed")
	ErrCacheCorruption = errors.New("cache corruption")
)

const (
	englishSlug        = "english"
	chineseSlug        = "chinese"
	japaneseSlug       = "japanese"
	englishCachePeriod = 3 * time.Minute
	// englishCachePeriod   = 3 * time.Millisecond
	chineseCachePeriod   = 5 * time.Minute
	japaneseCachePeriod  = 10 * time.Minute
	tutorInfoCachePeriod = 10 * time.Minute
	defaultAsyncDeadline = 30 * time.Second
)

var (
	langCacheKey  = "langslug:%s:cache"
	tutorCacheKey = "tutor:%s:cache"
)

func InitCacheLayer(storage repository.Repository, logger *zap.Logger) repository.Repository {
	return &cacheLayer{
		db:     storage,
		logger: logger,
		cache:  gocache.New(gocache.NoExpiration, gocache.NoExpiration),
		ml:     NewMultipleLock(),
	}
}

type cacheLayer struct {
	db     repository.Repository
	logger *zap.Logger
	cache  *gocache.Cache
	ml     MultipleLock
}

type cacheEntry struct {
	updatedAt time.Time
	data      interface{}
}

func getCachePeriod(langSlug string) time.Duration {
	switch langSlug {
	case englishSlug:
		return englishCachePeriod
	case chineseSlug:
		return chineseCachePeriod
	case japaneseSlug:
		return japaneseCachePeriod
	default:
		return englishCachePeriod
	}
}

func (c *cacheLayer) GetTutors(ctx context.Context, langSlug string) ([]repository.TutorInfo, error) {
	cacheKey := fmt.Sprintf(langCacheKey, langSlug)

	// cache hit
	rawdata, found := c.cache.Get(cacheKey)
	if found {
		entry, ok := rawdata.(cacheEntry)
		if !ok {
			return nil, ErrCacheCorruption
		}
		oriData, ok := entry.data.([]repository.TutorInfo)
		if !ok {
			return nil, ErrCacheCorruption
		}

		cachePeriod := getCachePeriod(langSlug)
		if entry.updatedAt.Add(cachePeriod).Before(time.Now()) {
			// invalid cache, trigger background job to update cache
			ctx, cancel := context.WithDeadline(ctx, time.Now().Add(defaultAsyncDeadline))
			go func() {
				defer cancel()
				tutors, err := c.db.GetTutors(ctx, langSlug)
				if err != nil {
					c.logger.Warn("GetTutors error", zap.Error(err))
					return
				}
				c.updateCache(cacheKey, tutors, time.Now())
			}()
		}
		c.logger.Info("[Cache hit] return tutors info", zap.String("langSlug", langSlug))
		return oriData, nil
	}

	// cache missed
	c.logger.Debug("[Cache miss] try to got lock", zap.String("langSlug", langSlug))
	c.ml.Lock(cacheKey)
	defer c.ml.Unlock(cacheKey)
	c.logger.Debug("[Cache miss] got lock", zap.String("langSlug", langSlug))

	// goroutine gets the lock, check the cache again to prevent cache stampede
	rawdata, found = c.cache.Get(cacheKey)
	if found {
		c.logger.Debug("[Cache hit] good! just return the cache value", zap.String("langSlug", langSlug))
		entry, ok := rawdata.(cacheEntry)
		if !ok {
			return nil, ErrCacheCorruption
		}
		oriData, ok := entry.data.([]repository.TutorInfo)
		if !ok {
			return nil, ErrCacheCorruption
		}
		return oriData, nil
	}

	// responsible to update cache
	tutors, err := c.db.GetTutors(ctx, langSlug)
	if err != nil {
		c.logger.Warn("GetTutors error", zap.Error(err))
		return nil, err
	}
	c.updateCache(cacheKey, tutors, time.Now())
	return tutors, nil
}

func (c *cacheLayer) GetTutor(ctx context.Context, tutorSlug string) (*repository.TutorInfo, error) {
	cacheKey := fmt.Sprintf(tutorCacheKey, tutorSlug)
	// cache hit
	rawdata, found := c.cache.Get(cacheKey)
	if found {
		entry, ok := rawdata.(cacheEntry)
		if !ok {
			return nil, ErrCacheCorruption
		}
		oriData, ok := entry.data.(repository.TutorInfo)
		if !ok {
			return nil, ErrCacheCorruption
		}

		if entry.updatedAt.Add(tutorInfoCachePeriod).Before(time.Now()) {
			// invalid cache, trigger background job to update cache
			ctx, cancel := context.WithDeadline(ctx, time.Now().Add(defaultAsyncDeadline))
			go func() {
				defer cancel()
				tutor, err := c.db.GetTutor(ctx, tutorSlug)
				if err != nil {
					c.logger.Warn("GetTutor error", zap.Error(err))
					return
				}
				c.logger.Debug("async update GetTutor", zap.String("tutorSlug", tutorSlug))
				c.updateCache(cacheKey, *tutor, time.Now())
			}()
		}
		c.logger.Info("[Cache hit] return tutor info", zap.String("tutorSlug", tutorSlug))
		return &oriData, nil
	}

	// cache missed
	c.logger.Debug("[Cache miss] try to got lock", zap.String("tutorSlug", tutorSlug))
	c.ml.Lock(cacheKey)
	defer c.ml.Unlock(cacheKey)
	c.logger.Debug("[Cache miss] got lock", zap.String("tutorSlug", tutorSlug))

	// goroutine gets the lock, check the cache again to prevent cache stampede
	rawdata, found = c.cache.Get(cacheKey)
	if found {
		c.logger.Debug("[Cache hit] good! just return the cache value", zap.String("tutorSlug", tutorSlug))
		entry, ok := rawdata.(cacheEntry)
		if !ok {
			return nil, ErrCacheCorruption
		}
		oriData, ok := entry.data.(repository.TutorInfo)
		if !ok {
			return nil, ErrCacheCorruption
		}
		return &oriData, nil
	}

	// responsible to update cache
	tutor, err := c.db.GetTutor(ctx, tutorSlug)
	if err != nil {
		c.logger.Warn("GetTutor error", zap.Error(err))
		return nil, err
	}
	c.updateCache(cacheKey, *tutor, time.Now())
	return tutor, nil
}

func (c *cacheLayer) updateCache(key string, data interface{}, updatedAt time.Time) {
	entry := cacheEntry{
		updatedAt: updatedAt,
		data:      data,
	}
	c.logger.Info("[Cache update] cache data", zap.String("key", key))
	c.cache.Set(key, entry, gocache.NoExpiration)
}
