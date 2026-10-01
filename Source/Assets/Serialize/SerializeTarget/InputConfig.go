package serializetarget

import (
	inputconfig "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/InputConfig"
	"fmt"

	input "github.com/quasilyte/ebitengine-input"
)

// inputconfig.InputConfig を構築するためのシリアライズ構造体
type InputConfig struct {
	KeyMap map[string][]string
}

func (ic InputConfig) Construct() (any, error) {
	keymap := input.Keymap{}

	for serikey, serivalue := range ic.KeyMap {
		key, err := inputconfig.ActionConst(serikey)
		if err != nil {
			return nil, fmt.Errorf("inputconfig.ActionConst() でエラーが発生しました。\n%w", err)
		}
		value := []input.Key{}
		for _, keystring := range serivalue {
			inputkey, err := input.ParseKey(keystring)
			if err != nil {
				return nil, fmt.Errorf("input.ParseKey(\"%s\")  でエラーが発生しました。\n%w", inputkey, err)
			}
			value = append(value, inputkey)
		}
		keymap[key] = value
	}

	return keymap, nil
}
