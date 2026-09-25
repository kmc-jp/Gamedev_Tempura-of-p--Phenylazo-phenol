package serialize

import (
	"fmt"
	"reflect"
)

type SerializeData struct {
	Name string
	Type reflect.Type
	// /Assets/Serialize からの相対パス
	Path string
}

func NewSerializeData(name string, Type reflect.Type, path string, serializeManager *SerializeManager) (*SerializeData, error) {
	cleanedPath, err := serializeManager.CheckAndCleanPath(path)
	if err != nil {
		return nil, fmt.Errorf("path が不正です。\n%w", err)
	}

	id := SerializeData{
		Name: name,
		Path: cleanedPath,
	}

	if !Type.Implements(reflect.TypeFor[Serializable]()) {
		return nil, fmt.Errorf("Type %s は Serializable interface を実装していません。", Type.Name())
	}

	return &id, nil
}
