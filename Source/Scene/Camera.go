package scene

import (
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"
	"math"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
)

// カメラの本実装
// 場所を指定したらその場所を中心として映す
type Camera struct {
	center     Utils.Position
	layers     map[string]int // 値のintは正の数
	layerCount int
}

type drawImage struct {
	image  *ebiten.Image
	option *ebiten.DrawImageOptions
}

// Draw implements [Drawer].
func (c Camera) Draw(screen *ebiten.Image, objects []Entity) {
	scSize := screen.Bounds().Size()
	images := make([][]drawImage, c.layerCount)
	for _, obj := range objects {
		image, Layer := obj.Draw()
		layernum := c.layers[Layer] // 存在しない場合はゼロ
		pos := obj.Position().Offset(c.center)

		posx, posy := pos.X+float64(scSize.X)/2, pos.Y+float64(scSize.Y)/2
		sizex, sizey := image.Bounds().Dx(), image.Bounds().Dy()
		iposx, iposy := math.Floor(posx-float64(sizex)/2), math.Floor(posy-float64(sizey)/2)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(iposx, iposy)

		// 各レイヤーに割り当て
		if images[layernum] == nil {
			images[layernum] = make([]drawImage, len(objects))
		}
		images[layernum] = append(images[layernum], drawImage{image: image, option: op})
	}

	for _, layer := range images {
		for _, draw := range layer {
			screen.DrawImage(draw.image, draw.option)
		}
	}
}

func (c *Camera) SetCenter(center Utils.Position) {
	c.center = center
}

var _ Drawer = &Camera{}

type kvp struct {
	key   string
	value int
}

func NewCamera(layers map[string]int) Utils.Factory[*Camera] {
	layers, err := compressValues(layers)
	if err != nil {
		return func() (*Camera, error) {
			return nil, fmt.Errorf("compressValues(layers) でエラーが発生しました。,%w", err)
		}
	}
	return func() (*Camera, error) {
		return &Camera{
			center:     Utils.Position{},
			layers:     layers,
			layerCount: len(layers),
		}, nil
	}
}

func compressValues(data map[string]int) (map[string]int, error) {
	// 1. 各要素をスライスに集める
	pairs := make([]kvp, 0, len(data))
	for k, v := range data {
		pairs = append(pairs, kvp{k, v})
	}

	// 2. Valueを基準に昇順ソートする
	slices.SortFunc(pairs, func(a, b kvp) int {
		return a.value - b.value
	})

	// 3. 事前チェック: 隣り合うValueを比較して重複があればエラーを返す
	// ※ソート済みのため、重複があれば必ず隣り合います
	for i := 1; i < len(pairs); i++ {
		if pairs[i].value == pairs[i-1].value {
			return nil, fmt.Errorf("重複するValueを検出しました: %d (キー: %s, %s)",
				pairs[i].value, pairs[i-1].key, pairs[i].key)
		}
	}

	// 4. 0 ~ n の連番に圧縮して新しいマップを作成
	compressed := make(map[string]int, len(data))
	for i, p := range pairs {
		compressed[p.key] = i
	}

	return compressed, nil
}
