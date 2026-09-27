package serialize

import (
	"fmt"
	"io"
	"io/fs"
	"reflect"

	"github.com/BurntSushi/toml"
)

type TomlLoader struct {
	cacheToml map[string]Serializable
}

func NewTomlLoader() (*TomlLoader, error) {
	tl := TomlLoader{
		cacheToml: map[string]Serializable{},
	}
	return &tl, nil
}

func (tl TomlLoader) Load(data SerializeData, fs fs.FS) (Serializable, error) {
	_, ok := tl.cacheToml[data.Name]
	if !ok {
		item, err := tl.fetch(data, fs)
		if err != nil {
			return nil, fmt.Errorf("TomlLoader.fetch() でエラーが発生しました。\n%w", err)
		}
		tl.cacheToml[data.Name] = item
	}
	item, _ := tl.cacheToml[data.Name]
	return item, nil
}

func (tl TomlLoader) fetch(data SerializeData, fs fs.FS) (Serializable, error) {
	// まずはファイルの中身を string で取り出す
	file, err := fs.Open(data.Path)
	if err != nil {
		return nil, fmt.Errorf("fs.Open(data.Path) (Name : %s, Path : %s) でエラーが発生しました。\n%w", data.Name, data.Path, err)
	} else {
		defer file.Close()
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll(file) でエラーが発生しました。\n%w", err)
	}

	tomldata := string(bytes)

	// デシリアライズ実行
	valPtr := reflect.New(data.Type).Interface() // ここにデシリアライズされたデータが入る
	if _, err := toml.Decode(tomldata, valPtr); err != nil {
		return nil, fmt.Errorf("toml.Decode(tomldata, valPtr) でエラーが発生しました。\n%w", err)
	}

	var deserialized Serializable = reflect.ValueOf(valPtr).Elem().Interface().(Serializable) // ポインタを Serializable 型に変換

	return deserialized, nil
}
