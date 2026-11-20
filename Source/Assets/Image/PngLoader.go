package imageassets

import (
	assetsutil "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets/AssetsUtils"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// LoadPNG はAseprite用のjsonを持たない単体のPNGを読み込む。
// pngPath は /Assets/Images からの相対パス。同じパスは読み込み結果を使い回す。
func (im *ImageManager) LoadPNG(pngPath string) (*ebiten.Image, error) {
	cleaned, err := assetsutil.CheckAndCleanPath(pngPath, im.imageFS)
	if err != nil {
		return nil, fmt.Errorf("pngPath が不正です。\n%w", err)
	}
	if img, ok := im.pngCache[cleaned]; ok {
		return img, nil
	}
	img, _, err := ebitenutil.NewImageFromFileSystem(im.imageFS, cleaned)
	if err != nil {
		return nil, fmt.Errorf("ebitenutil.NewImageFromFileSystem() でエラーが発生しました。\npath : %s\n%w", cleaned, err)
	}
	im.pngCache[cleaned] = img
	return img, nil
}
