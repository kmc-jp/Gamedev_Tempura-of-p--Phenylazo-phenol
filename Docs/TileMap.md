# タイルマップ

## tilemap.TileData

タイルの初期配置情報
要はどこに何のタイルがあるか
- map[TilePosition]Tile

## tilemap.TilePosition

{
    x int
    y int
}

## Component.TileMap

タイルの配置位置に関する情報
初期配置についても担当

### New(mapdata tilemap.TileData)
- Tileを大量生成する

## Component.Tile

一枚一枚のタイルに関する情報
動かない

## Component.TileWalker

タイルマップ上を動くものに関する情報
動く
