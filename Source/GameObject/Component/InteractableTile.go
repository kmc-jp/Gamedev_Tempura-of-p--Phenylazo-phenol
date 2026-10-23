package component

import (
	gameobject "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/GameObject"
	tilemap "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/TileMap"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"
)

type InteractableTile interface {
	gameobject.Component
	Interact() error
}

// tilemap.InteractType を、対応するコンポーネントのファクトリに変換する
// nil かもしれないので対応できるようにせよ
func InteractiveTiles(Type tilemap.InteractType) (Utils.Factory[gameobject.Component], error) {
	switch Type {
	case tilemap.None:
		return nil, nil
	case tilemap.Shelf:
		return Utils.CastFactory[*ShelfTile, gameobject.Component](NewShelfTile()), nil
	default:
		return nil, fmt.Errorf("未定義の Type が渡されました。\nInteractiveTiles() または tilemap.InteractType の定義を確認してください。")
	}
}
