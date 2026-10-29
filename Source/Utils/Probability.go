package Utils

import (
	"fmt"
	"math/rand"
)

// 0 < p < 1 の確率を入力
func WithProbability(p float64) bool {
	r := rand.Float64()
	if r <= p {
		return true
	} else {
		return false
	}
}

// 0 < p < 1 以外の場合にerror返却
func WithProbability_strict(p float64) (bool, error) {
	if p >= 1 {
		return false, fmt.Errorf("Percent_strict()に1以上の値が入力されました")
	}
	if p <= 0 {
		return false, fmt.Errorf("Percent_strict()に0以下の値が入力されました")
	}
	r := rand.Float64()
	if r <= p {
		return true, nil
	} else {
		return false, nil
	}
}
