package Utils

import (
	"fmt"
	"math/rand"
)

func Percent(p float64) bool {
	r := rand.Float64()
	if r <= p {
		return true
	} else {
		return false
	}
}

func Percent_strict(p float64) (bool, error) {
	if p >= 1 {
		return false, fmt.Errorf("Percent_strict()に1以上の値が入力されました")
	}
	r := rand.Float64()
	if r <= p {
		return true, nil
	} else {
		return false, nil
	}
}
