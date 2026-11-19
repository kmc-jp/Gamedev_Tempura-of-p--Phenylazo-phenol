package component

import (
	gameobject "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/GameObject"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

// 複数のレンダリングを重ねるレンダー。
// レイヤーはmainのものが返る。
// subrendering と subupdate の一方がfalseならSubは描画停止
// subrendering は false だが subupdate は true のとき、アニメーションだけ裏で進行する
type MultiRender struct {
	main         gameobject.Render
	subs         []gameobject.Render
	subrendering bool
	subupdate    bool

	imgcache      *ebiten.Image
	imgqueuecache []*ebiten.Image
	dio           *ebiten.DrawImageOptions
	dioqueuecache []*ebiten.DrawImageOptions
}

func NewMultiRender(main gameobject.Render, subs []gameobject.Render) *MultiRender {
	mr := MultiRender{
		main:          main,
		subs:          subs,
		subrendering:  false,
		subupdate:     false,
		imgcache:      ebiten.NewImage(1, 1),
		imgqueuecache: make([]*ebiten.Image, len(subs)+1),
		dioqueuecache: make([]*ebiten.DrawImageOptions, len(subs)+1),
		dio:           &ebiten.DrawImageOptions{},
	}

	return &mr
}

// Factory 状態で受け取って Factory を返す
func MultiRenderDecorator(subFactories ...Utils.Factory[gameobject.Render]) func(Utils.Factory[gameobject.Render]) Utils.Factory[gameobject.Render] {
	return func(mainFactory Utils.Factory[gameobject.Render]) Utils.Factory[gameobject.Render] {
		return func() (gameobject.Render, error) {
			main, err := mainFactory()
			if err != nil {
				return nil, fmt.Errorf("MultiRender の mainFactory での生成に失敗しました。\n%w", err)
			}

			subs := make([]gameobject.Render, len(subFactories))
			for i, f := range subFactories {
				sub, err := f()
				if err != nil {
					return nil, fmt.Errorf("MultiRender の subFactories[%d] での生成に失敗しました。\n%w", i, err)
				}
				subs[i] = sub
			}

			return NewMultiRender(main, subs), nil
		}
	}
}

func (mr *MultiRender) Draw() (image *ebiten.Image, DIO *ebiten.DrawImageOptions, Layer string) {
	if mr.subupdate && mr.subrendering {
		// mr.imgcache.Clear()
		mr.imgqueuecache[0], mr.dioqueuecache[0], Layer = mr.main.Draw()
		mainimg := mr.imgqueuecache[0]
		// mr.imgcache.DrawImage(img, dio)

		for i, rnd := range mr.subs {
			mr.imgqueuecache[i+1], mr.dioqueuecache[i+1], _ = rnd.Draw()
			// mr.imgcache.DrawImage(img, dio)
		}

		if mr.imgcache.Bounds().Dx() != mainimg.Bounds().Dx() || mr.imgcache.Bounds().Dy() != mainimg.Bounds().Dy() {
			mr.imgcache = ebiten.NewImage(mainimg.Bounds().Dx(), mainimg.Bounds().Dy())
		} else {
			mr.imgcache.Clear()
		}

		for i := 0; i < len(mr.imgqueuecache); i++ {
			mr.imgcache.DrawImage(mr.imgqueuecache[i], mr.dioqueuecache[i])
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
