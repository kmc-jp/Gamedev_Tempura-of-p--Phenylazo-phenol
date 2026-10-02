package inputconfig

import (
	"fmt"

	input "github.com/quasilyte/ebitengine-input"
)

const (
	ActionUnknown input.Action = iota
	ActionMoveLeft
	ActionMoveRight
	ActionMoveUp
	ActionMoveDown
)

func ActionConst(name string) (input.Action, error) {
	switch name {
	case "ActionUnknown":
		return ActionUnknown, nil
	case "ActionMoveLeft":
		return ActionMoveLeft, nil
	case "ActionMoveRight":
		return ActionMoveRight, nil
	case "ActionMoveUp":
		return ActionMoveUp, nil
	case "ActionMoveDown":
		return ActionMoveDown, nil
	default:
		return ActionUnknown, fmt.Errorf("%s というアクションは見つかりませんでした。", name)
	}
}

type InputConfig struct {
	keymap input.Keymap // map[input.Action]([]Key)
}

func NewInputConfig(keymap input.Keymap) Utils.Factory[*InputConfig] {
	return func() (*InputConfig, error) {
		return &InputConfig{keymap: keymap}, nil
	}
}

func (ic InputConfig) Setup() (*input.Handler, error) {
	// inputSystem を初期化
	var inputSystem input.System
	inputSystem.Init(input.SystemConfig{
		DevicesEnabled: input.AnyDevice, // キーボード・マウス・ゲームパッドなどを有効化
	})

	// Player の InputSystem
	inputHandler := inputSystem.NewHandler(0, ic.keymap)

	return inputHandler, nil
}
