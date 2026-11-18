package component

import (
	gameobject "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/GameObject"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

// 複数のレンダリングを重ねるレンダー。
// レイヤーはmainのものが返る。
type MultiRender struct {
	main         gameobject.Render
	subs         []gameobject.Render
	subrendering bool
	subupdate    bool

	imgcache *ebiten.Image
	dio      *ebiten.DrawImageOptions
}

func (mr *MultiRender) Draw() (image *ebiten.Image, DIO *ebiten.DrawImageOptions, Layer string) {
	if mr.subupdate && mr.subrendering {
		mr.imgcache.Clear()
		img, dio, Layer := mr.main.Draw()
		mr.imgcache.DrawImage(img, dio)

		for _, rnd := range mr.subs {
			img, dio, _ := rnd.Draw()
			mr.imgcache.DrawImage(img, dio)
		}

		return mr.imgcache, mr.dio, Layer
	} else {
		return mr.main.Draw()
	}
}

func (mr *MultiRender) Update() error {
	if err := mr.main.Update(); err != nil {
		return fmt.Errorf("MultiRender.main.Update() でエラーが発生しました。\n%w", err)
	}

	if mr.subupdate {
		for i, rnd := range mr.subs {
			if err := rnd.Update(); err != nil {
				return fmt.Errorf("MultiRender.subs[%d].Update() でエラーが発生しました。\n%w", i, err)
			}
		}
	}

	return nil
}

func (mr *MultiRender) EnableSubRendering(active bool) {
	mr.subrendering = active
}

func (mr *MultiRender) EnableSubUpdate(active bool) {
	mr.subupdate = active
}

var _ gameobject.Render = (*MultiRender)(nil)
