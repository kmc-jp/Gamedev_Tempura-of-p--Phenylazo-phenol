# ぱらふぇにるあぞふぇのーるのてんぷら

## ゲーム企画

ジャンク屋でPCパーツを発掘するゲーム 

### システム

- ローグライト
- 収集
- 2Dステージ
- マルチステージ制(予定)

### フレーバー

- ドット絵
- ジャンクPC
- 綺麗～小汚いハード屋

## For Engineer
### 言語・パッケージ

- Golang
- ebitengine : ゲームエンジン
- asepriteをどうにかして読むパッケージ

### 環境構築

#### Go
https://go.dev/ \
または
```Powershell
winget install GoLang.Go
```

#### ebitengine
https://ebitengine.org/en/documents/install.html?os=windows


### フォルダ構成

ディレクトリ追加したら書いておきましょう \
全ファイル書く必要はないけど何してるディレクトリか分かるようにすべし

```
/
├── Assets
│   ├── Images
│   │   └── img
│   └── Sounds
│       └── snd
├── Source
│   ├── HogeScene
│   │   └── hoge
│   ├── Utils
│   └── Src
├── go.mod
├── go.sum
├── LISENCE
├── main.go
└── README.md
```

### 命名規則等

#### コミット

- 1行目　タイトル
- 2行目以降　変えた理由 (必要なら変えなかった理由)

#### ブランチ名

- 英語
- mainから直接生やす

#### PR

- タイトルだけでわかるようにしましょう
- 対応するIssueを書く
    - 説明欄に「 close #(Issue番号)」で対応するものが閉じられます \
    (mainにマージしたときのみ)

#### Issue

- タイトル
    - 実装すべきこと
    - 相談のタイトル
- 内容
    - 中身

### コード規則

- goに入ればgoに従え
- commit前に見た目を整える
- go fmt hoge.goで形は整う
  - (というかvscodeの拡張機能で、ctrl+s時にたぶん自動で走らせてる)
- err は最上位の関数まで返して戻り値にする


### タスク管理等

- Github の Issue を積極的に使っていきたい
- 更新のたびにBotが知らせてくれる用のチャンネルをつくる