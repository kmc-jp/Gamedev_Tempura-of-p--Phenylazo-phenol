package component

import (
	gameobject "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/GameObject"
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"
	tilemap "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/TileMap"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"
)

type InteractableTile interface {
	gameobject.Component
	Interact(*Player) error
}

// tilemap.InteractType を、対応するコンポーネントのファクトリに変換する
// nil かもしれないので対応できるようにせよ
func InteractiveTiles(Type tilemap.InteractType, scene *scene.PlayScene, render gameobject.Render) (Utils.Factory[InteractableTile], error) {
	switch Type {
	case tilemap.None:
		return nil, nil
	case tilemap.Shelf:
		mr, ok := render.(*MultiRender)
		if !ok {
			return nil, fmt.Errorf("*MultiRenderではないRenderがついたShelfがあります。")
		}
		return Utils.CastFactory[*ShelfTile, InteractableTile](NewShelfTile(scene, mr)), nil
	case tilemap.Register:
		return Utils.CastFactory[*RegisterTile, InteractableTile](NewRegisterTile(scene)), nil
	default:
		return nil, fmt.Errorf("未定義の Type が渡されました。\nInteractiveTiles() または tilemap.InteractType の定義を確認してください。")
	}
}
