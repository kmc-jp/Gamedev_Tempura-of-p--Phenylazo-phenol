package construct

import (
	assets "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets"
	gameobject "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/GameObject"
	component "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/GameObject/Component"
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"
	tilemap "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/TileMap"
	"Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/goaseprite"
)

// 個々の GameObject の定義を書く

// タイルマップ関連
func NewTileObject(name string, data tilemap.TileData, scene *scene.PlayScene, NewTransform Utils.Factory[gameobject.Transform], NewRender Utils.Factory[gameobject.Render]) Utils.Factory[*gameobject.GameObject] {
	return func() (*gameobject.GameObject, error) {
		o, err := gameobject.NewGameObject(name, scene, NewTransform, NewRender)()
		if err != nil {
			return o, fmt.Errorf("newStdObject() でエラーが発生しました。\n%w", err)
		}
		t, err := component.NewTile(data)()
		if err != nil {
			return o, fmt.Errorf("NewTile() でエラーが発生しました。\n%w", err)
		}
		err = o.AddComponent(t)
		if err != nil {
			return o, fmt.Errorf("GameObject.AddComponent() でエラーが発生しました。\n%w", err)
		}
		return o, nil
	}
}

func NewTileMapObject(
	name string,
	mapdata tilemap.TileMapData,
	PlScene *scene.PlayScene,
	NewTransform func(Utils.Position) Utils.Factory[gameobject.Transform],
	NewRender func(asp *goaseprite.File, img *ebiten.Image, layer string) Utils.Factory[gameobject.Render],
	assetsManager *assets.AssetsManager,
) Utils.Factory[*gameobject.GameObject] {
	return func() (*gameobject.GameObject, error) {
		newEmptyRender := func() (gameobject.Render, error) {
			r, e := component.NewEmptyRender()
			return r, e
		}
		o, err := gameobject.NewGameObject(
			name,
			PlScene,
			NewTransform(Utils.NewZeroVec2().VtoP()),
			newEmptyRender,
		)()
		if err != nil {
			return o, fmt.Errorf("newStdObject() でエラーが発生しました。\n%w", err)
		}

		tm, err := component.NewTileMap(mapdata)()
		if err != nil {
			return o, fmt.Errorf("NewTileMap() でエラーが発生しました。\n%w", err)
		}
		err = o.AddComponent(tm)
		if err != nil {
			return o, fmt.Errorf("GameObject.AddComponent() でエラーが発生しました。\n%w", err)
		}

		// タイル配置
		newTileObj := func(
			name string,
			data tilemap.TileData,
			Scene *scene.PlayScene,
			layer string,
			transformFactry Utils.Factory[gameobject.Transform],
		) Utils.Factory[*gameobject.GameObject] {
			return func() (*gameobject.GameObject, error) {
				asp, img, err := assetsManager.Image.Load(data.ImageData)
				if err != nil {
					mock, _ := gameobject.NewMockGameObject(name)()
					return mock, fmt.Errorf("assetsManager.Image.Load(data.ImageData) でエラーが発生しました。\n%w", err)
				}
				toFactory := NewTileObject(name, data, Scene, transformFactry, NewRender(asp, img, layer))
				return toFactory()
			}
		}

		err = tm.SetGameObject(PlScene, newTileObj, NewTransform)
		if err != nil {
			return o, fmt.Errorf("TileMap.SetGameObject() でエラーが発生しました。\n%w", err)
		}

		return o, nil
	}
}

func NewCharactor(
	name string,
	position tilemap.TilePosition,
	TileData tilemap.TileData,
	TileMap *component.TileMap,
	PlScene *scene.PlayScene,
	NewTransform func(Utils.Position) Utils.Factory[gameobject.Transform],
	NewRender func(asp *goaseprite.File, img *ebiten.Image, layer string) Utils.Factory[gameobject.Render],
	assetsManager *assets.AssetsManager,
) Utils.Factory[*gameobject.GameObject] {
	return func() (*gameobject.GameObject, error) {
		// キャラクター画像設定
		asp, img, err := assetsManager.Image.Load(TileData.ImageData)
		if err != nil {
			return nil, fmt.Errorf("assetsManager.Image.Load(TileData.ImageData) でエラーが発生しました。\n%w", err)
		}

		o, err := gameobject.NewGameObject(name, PlScene, NewTransform(position.Position.ToVec2().VtoP()), NewRender(asp, img, TileData.Layer))()
		if err != nil {
			return o, fmt.Errorf("newStdObject() でエラーが発生しました。\n%w", err)
		}

		// TileWalker
		TW, err := component.NewTileWalker(TileData, position, o.Transform(), TileMap, o.Render().(*component.AnimRender))()
		if err != nil {
			return nil, fmt.Errorf("component.NewTileWalker(TileData, position)() でエラーが発生しました。\n%w", err)
		}
		err = o.AddComponent(TW)
		if err != nil {
			return nil, fmt.Errorf("o.AddComponent(TW) でエラーが発生しました。\n%w", err)
		}

		return o, nil
	}
}

func NewPlayer(
	name string,
	position tilemap.TilePosition,
	TileData tilemap.TileData,
	TileMap *component.TileMap,
	PlScene *scene.PlayScene,
	NewTransform func(Utils.Position) Utils.Factory[gameobject.Transform],
	NewRender func(asp *goaseprite.File, img *ebiten.Image, layer string) Utils.Factory[gameobject.Render],
	assetsManager *assets.AssetsManager,
) Utils.Factory[*gameobject.GameObject] {
	return func() (*gameobject.GameObject, error) {
		plobj, err := NewCharactor(name, position, TileData, TileMap, PlScene, NewTransform, NewRender, assetsManager)()
		if err != nil {
			return nil, fmt.Errorf("NewCharactor() でエラーが発生しました。\n%w", err)
		}

		// player (Component)
		plc, err := component.NewPlayer(PlScene, plobj, plobj.Render().(*component.AnimRender))()
		if err != nil {
			return nil, fmt.Errorf("component.NewPlayer(PlScene, pl)() でエラーが発生しました。\n%w", err)
		}
		err = plobj.AddComponent(plc)
		if err != nil {
			return nil, fmt.Errorf("pl.AddComponent(plc) でエラーが発生しました。\n%w", err)
		}

		return plobj, nil
	}
}
