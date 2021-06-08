package main

import (
	"cachepractice/repository"
	"encoding/json"
	"fmt"
)

func main() {
	data := repository.TutorInfo{
		TutorID:       "xxx",
		TutorSlug:     "foobar",
		TutorName:     "Amazing Teacher 1",
		TutorHeadline: "Hi I'm a English Teacher",
		PriceInfo: repository.PriceInfo{
			Trial:  5,
			Normal: 10,
		},
		TeachingLangs: []int{123, 231},
	}
	bytes, err := json.MarshalIndent(data, "", "    ")
	fmt.Println(err)
	fmt.Println(string(bytes))
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
