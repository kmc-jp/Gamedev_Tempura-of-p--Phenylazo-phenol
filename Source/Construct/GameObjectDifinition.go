package construct

import (
	gameobject "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/GameObject"
	component "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/GameObject/Component"
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"
	tilemap "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/TileMap"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"
)

// 個々の GameObject の定義を書く

// タイルマップ関連
func NewTileObject(name string, data tilemap.TileData, scene *scene.PlayScene, NewTransform Utils.Factory[gameobject.Transform], NewRender Utils.Factory[gameobject.Render]) Utils.Factory[gameobject.GameObject] {
	o, err := gameobject.NewGameObject(name, scene, NewTransform, NewRender)()
	if err != nil {
		return func() (*gameobject.GameObject, error) {
			return o, fmt.Errorf("newStdObject() でエラーが発生しました。\n%w", err)
		}
	}
	t, err := component.NewTile(data)()
	if err != nil {
		return func() (*gameobject.GameObject, error) {
			return o, fmt.Errorf("NewTile() でエラーが発生しました。\n%w", err)
		}
	}
	o.AddComponent(t)
	return func() (*gameobject.GameObject, error) {
		return o, nil
	}
}

// test
var _ Utils.Factory[gameobject.GameObject] = NewTileObject("", tilemap.TileData{}, nil, nil, nil)
