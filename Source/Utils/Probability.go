package Utils

import (
	"fmt"
	"math/rand"
)

// 確率pでtrueを返す
func Percent(p float64) bool {
	r := rand.Float64()
	if r <= p {
		return true
	} else {
		return false
	}
}

// 確率pでtrueを返すがpに1以上の値が入力された場合エラーを返す
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
