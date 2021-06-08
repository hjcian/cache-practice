//
// this implementation is just for testing purposes
//
package repository

import (
	"context"
)

func InitFakeDB() (Repository, error) {
	return &fakedb{}, nil
}

type fakedb struct{}

func (f *fakedb) GetTutors(ctx context.Context, langSlug string) ([]TutorInfo, error) {
	switch langSlug {
	case "english":
		return []TutorInfo{englishTutor1, englishTutor2}, nil
	case "chinese":
		return []TutorInfo{chineseTutor1}, nil
	case "japanese":
		return []TutorInfo{japaneseTutor1}, nil
	}
	return nil, ErrRecordNotFound
}

func (f *fakedb) GetTutor(ctx context.Context, tutorSlug string) (*TutorInfo, error) {
	fakeTutors := getFakeTutors()
	tutorInfo, ok := fakeTutors[tutorSlug]
	if !ok {
		return nil, ErrRecordNotFound
	}
	return &tutorInfo, nil
}

var fakeTutors map[string]TutorInfo

func getFakeTutors() map[string]TutorInfo {
	if fakeTutors != nil {
		return fakeTutors
	}
	fakeTutors = make(map[string]TutorInfo)
	fakeTutors[englishTutor1.TutorSlug] = englishTutor1
	fakeTutors[englishTutor2.TutorSlug] = englishTutor2
	fakeTutors[chineseTutor1.TutorSlug] = chineseTutor1
	fakeTutors[japaneseTutor1.TutorSlug] = japaneseTutor1
	return fakeTutors
}

var (
	englishTutor1 = TutorInfo{
		TutorID:       "1",
		TutorSlug:     "foo-bar",
		TutorName:     "Amazing Teacher 1",
		TutorHeadline: "Hi I'm a English Teacher",
		PriceInfo: PriceInfo{
			Trial:  5,
			Normal: 10,
		},
		TeachingLangs: []int{123, 321},
	}
	englishTutor2 = TutorInfo{
		TutorID:       "2",
		TutorSlug:     "bar-foo",
		TutorName:     "Amazing Teacher 2",
		TutorHeadline: "Hi I'm a English Teacher",
		PriceInfo: PriceInfo{
			Trial:  10,
			Normal: 20,
		},
		TeachingLangs: []int{123, 321},
	}
	chineseTutor1 = TutorInfo{
		TutorID:       "3",
		TutorSlug:     "chinese-foo-bar",
		TutorName:     "Amazing chinese Teacher 1",
		TutorHeadline: "Hi I'm a Chinese Teacher",
		PriceInfo: PriceInfo{
			Trial:  10,
			Normal: 20,
		},
		TeachingLangs: []int{123, 321},
	}
	japaneseTutor1 = TutorInfo{
		TutorID:       "4",
		TutorSlug:     "japanese-foo-bar",
		TutorName:     "Amazing japanese Teacher 1",
		TutorHeadline: "Hi I'm a japanese Teacher",
		PriceInfo: PriceInfo{
			Trial:  15,
			Normal: 25,
		},
		TeachingLangs: []int{123, 456},
	}
)
