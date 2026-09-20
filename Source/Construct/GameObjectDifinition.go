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
func NewTileObject(name string, data tilemap.TileData, scene *scene.PlayScene, NewTransform Utils.Factory[gameobject.Transform], NewRender Utils.Factory[gameobject.Render]) Utils.Factory[*gameobject.GameObject] {
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
	err = o.AddComponent(t)
	if err != nil {
		return func() (*gameobject.GameObject, error) {
			return o, fmt.Errorf("GameObject.AddComponent() でエラーが発生しました。\n%w", err)
		}
	}
	return func() (*gameobject.GameObject, error) {
		return o, nil
	}
}

func NewTileMapObject(name string, mapdata tilemap.TileMapData, PlScene *scene.PlayScene, NewTransform func(Utils.Position) Utils.Factory[gameobject.Transform], NewRender Utils.Factory[gameobject.Render]) Utils.Factory[*gameobject.GameObject] {
	o, err := gameobject.NewGameObject(name, PlScene, NewTransform(Utils.NewZeroVec2().VtoP()), NewRender)()
	if err != nil {
		return func() (*gameobject.GameObject, error) {
			return o, fmt.Errorf("newStdObject() でエラーが発生しました。\n%w", err)
		}
	}

	tm, err := component.NewTileMap(mapdata)()
	if err != nil {
		return func() (*gameobject.GameObject, error) {
			return o, fmt.Errorf("NewTileMap() でエラーが発生しました。\n%w", err)
		}
	}
	err = o.AddComponent(tm)
	if err != nil {
		return func() (*gameobject.GameObject, error) {
			return o, fmt.Errorf("GameObject.AddComponent() でエラーが発生しました。\n%w", err)
		}
	}

	// タイル配置
	newTileObj := func(
		name string,
		data tilemap.TileData,
		Scene *scene.PlayScene,
		transformFactry Utils.Factory[gameobject.Transform],
	) Utils.Factory[*gameobject.GameObject] {
		return func() (*gameobject.GameObject, error) {
			toFactory := NewTileObject(name, data, Scene, transformFactry, NewRender)
			return toFactory()
		}
	}

	err = tm.SetGameObject(PlScene, newTileObj, NewTransform)
	if err != nil {
		return func() (*gameobject.GameObject, error) {
			return o, fmt.Errorf("TileMap.SetGameObject() でエラーが発生しました。\n%w", err)
		}
	}

	return func() (*gameobject.GameObject, error) {
		return o, nil
	}
}
