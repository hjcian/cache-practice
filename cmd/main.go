package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// data := repository.TutorInfo{
	// 	TutorID:       "xxx",
	// 	TutorSlug:     "foobar",
	// 	TutorName:     "Amazing Teacher 1",
	// 	TutorHeadline: "Hi I'm a English Teacher",
	// 	PriceInfo: repository.PriceInfo{
	// 		Trial:  5,
	// 		Normal: 10,
	// 	},
	// 	TeachingLangs: []int{123, 231},
	// }
	// bytes, err := json.MarshalIndent(data, "", "    ")
	// fmt.Println(err)
	// fmt.Println(string(bytes))
	fmt.Println(time.Now().Add(10 * time.Minute).Before(time.Now()))
}

type CacheLayer struct {
	triggerChan chan struct{}
	mutex       sync.Mutex
}

func (c *CacheLayer) poc() {
	select {
	case <-c.getLock():
		fmt.Println("...and sleep")
		time.Sleep(2 * time.Second)
	default:
		fmt.Println("no lock, sleep one second")
		time.Sleep(time.Second)
	}
}

func (c *CacheLayer) getLock() (ok chan struct{}) {
	c.mutex.Lock()
	fmt.Println("Got lock and sleep")
	ok <- struct{}{}
	c.mutex.Unlock()
	return
}

func (c *CacheLayer) updateCache(wg *sync.WaitGroup, lang string) {
	fmt.Println("hi!", lang)
	select {
	case c.triggerChan <- struct{}{}:
		//
		fmt.Println(lang, "try to get data...")
		time.Sleep(time.Second)
		fmt.Println(lang, "done.")
		//
		<-c.triggerChan
		wg.Done()
	default:
		// already trigger cache update, other request just use old value
		wg.Done()
	}
}

// "data": {
//     "id": "xxx",
//     "slug": "foo-bar",
//     "name": "Amazing Teacher 1",
//     "headline": "Hi I'm a English Teacher",
//     "introduction": ".........",
//     "price_info": {
//       "trial": 5,
//       "normal": 10
//     },
//     "teaching_languages": [123,121]
//   }
