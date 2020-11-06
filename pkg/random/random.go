package random

import (
	"math/rand"
	"time"
)

// GetRandomNumberInRange returns a random number in range [min,max]
func GetRandomNumberInRange(min,max int) int {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	result:=r1.Intn(max-min+1)+min
	return result
}