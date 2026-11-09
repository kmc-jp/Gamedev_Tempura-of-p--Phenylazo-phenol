package ui

import (
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"
	utils "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type AnimatableUI struct {
	*RichUI

	animationManager *utils.PathAnimateManager
	textOffset       utils.Vec2
	animating        bool
}

// RichUI生成用の値が必要
func NewAnimatableUI(
	position utils.Position,
	size image.Point,
	background *ebiten.Image,
	text *Text,
) (*AnimatableUI, error) {
	richUI, err := NewRichUI(position, size, background, text)
	if err != nil {
		return nil, fmt.Errorf("NewRichUI() でエラーが発生しました: %w", err)
	}

	textOffset := utils.Vec2{}
	if text != nil {
		textOffset = text.Position().Offset(position)
	}

	return &AnimatableUI{
		RichUI:           richUI,
		animationManager: utils.NewPathAnimateManager(),
		textOffset:       textOffset,
	}, nil
}

func (a *AnimatableUI) Update(active bool) error {
	if err := a.RichUI.Update(active); err != nil {
		return fmt.Errorf("RichUI の更新に失敗しました: %w", err)
	}
	if !active || !a.animating {
		return nil
	}

	if a.RichUI.position.PathAnimate(a.animationManager) == utils.Move_Finished {
		a.animating = false
	}
	a.syncTextPosition()
	return nil
}

func (a *AnimatableUI) SetPosition(position utils.Position) {
	a.RichUI.SetPosition(position)
	a.syncTextPosition()
}

// to: 目標位置, fromspeed: 初速, time: 終了までの設定時間(0.1s); 数値の組み合わせと実際の挙動の整合性に注意
func (a *AnimatableUI) SetPathAnimation(
	animationType utils.Motiontype,
	to utils.Vec2,
	fromSpeed float64,
	time uint,
) error {
	if err := a.animationManager.SetPathAnimation(
		animationType,
		to,
		fromSpeed,
		time,
		&a.RichUI.position,
	); err != nil {
		return fmt.Errorf("UIアニメーションの設定に失敗しました: %w", err)
	}
	a.animating = true
	return nil
}

func (a *AnimatableUI) AnimationManager() *utils.PathAnimateManager {
	return a.animationManager
}

func (a *AnimatableUI) syncTextPosition() {
	if a.RichUI.text == nil {
		return
	}
	a.RichUI.text.SetPosition(utils.Position{
		X: a.RichUI.position.X + a.textOffset.X,
		Y: a.RichUI.position.Y + a.textOffset.Y,
	})
}

var _ scene.Element = (*AnimatableUI)(nil)
var _ scene.PositionElement = (*AnimatableUI)(nil)
