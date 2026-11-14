package serialize

import (
	"fmt"
	"io/fs"
)

type SerializeManager struct {
	tomlLoader  *TomlLoader
	serializeFS fs.FS
}

func NewSerializeManager(serializeFS fs.FS) (*SerializeManager, error) {
	tl, err := NewTomlLoader()
	if err != nil {
		return &SerializeManager{}, fmt.Errorf("NewTomlLoader() でエラーが発生しました。\n%w", err)
	}
	return &SerializeManager{tomlLoader: tl, serializeFS: serializeFS}, nil
}

func (sm SerializeManager) Load(data SerializeData) (any, error) {
	item, err := sm.tomlLoader.Load(data, sm.serializeFS)
	if err != nil {
		return nil, fmt.Errorf("tomlLoader.Load() でエラーが発生しました。\n%w", err)
	}
	return item, nil
}

func (sm SerializeManager) GetFS() fs.FS {
	return sm.serializeFS
}
