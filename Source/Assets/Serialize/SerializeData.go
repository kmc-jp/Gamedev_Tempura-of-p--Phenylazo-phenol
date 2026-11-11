package serialize

import (
	assetsutil "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Assets/AssetsUtils"
	"fmt"
	"reflect"
)

type SerializeData struct {
	Name string
	Type reflect.Type
	// /Assets/Serialize からの相対パス
	Path string
}

func NewSerializeData(name string, Type reflect.Type, path string, SerializeManager *SerializeManager) (*SerializeData, error) {
	cleanedPath, err := assetsutil.CheckAndCleanPath(path, SerializeManager.serializeFS)
	if err != nil {
		return nil, fmt.Errorf("path が不正です。\n%w", err)
	}

	if !Type.Implements(reflect.TypeFor[Serializable]()) {
		return nil, fmt.Errorf("Type %s は Serializable interface を実装していません。", Type.Name())
	}

	id := SerializeData{
		Name: name,
		Path: cleanedPath,
		Type: Type,
	}

	return &id, nil
}
