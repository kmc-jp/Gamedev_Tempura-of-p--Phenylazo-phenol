package ui

import (
	"os"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

func LoadFont(path string) (*text.GoTextFaceSource, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	src, err := text.NewGoTextFaceSource(f)
	if err != nil {
		return nil, err
	}
	return src, nil
}

// this func can be used as font option.
func FontFace(src *text.GoTextFaceSource, size float64) *text.GoTextFace {
	return &text.GoTextFace{Source: src, Size: size}
}
