# タイルマップ

## tilemap.TileMapData

タイルの初期配置情報
要はどこに何のタイルがあるか
- map[TilePosition]Tile

## tilemap.TilePosition

{
    // 左上の座標
    Position Utils.Vec2Int
    Layer    string
    // 標準座標からずらしたい場合
    Offset   Utils.Vec2Int
}

## tilemap.TileData

タイルの種類に関する情報を保持
複数タイル使う場合もここに書いておく

## Component.TileMap

タイルの配置位置に関する情報
初期配置・進入可否についても担当

### New(mapdata tilemap.TileData)
- Tileを大量生成する

## Component.Tile

一枚一枚のタイルに関する情報
動かない

## Component.TileWalker

タイルマップ上を動くものに関する情報
動く

## 実装していないものの一覧

- 複数タイル使う場合の処理(仮素材で動作確認済み、本番は素材待ち)
- 壁貼り(とりあえず素材待ち)
