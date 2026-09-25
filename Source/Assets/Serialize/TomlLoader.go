package serialize

import (
	"fmt"
	"io/fs"
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
	return nil, fmt.Errorf("TomlLoader.fetch() は未実装です。")
}
