package scenecanvas

import (
	canvas "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Canvas"
	game "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Game"
	scene "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene"
	stageinfo "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Scene/StageInfo"
	ui "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/UI"
	utils "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	color "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils/Color"
	"fmt"
	"maps"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
)

// SelectStageCanvas は SelectStageScene のUI構成である。路線図と操作案内を持つ。
// 選択中の駅と現在地はSceneが持つため、Sceneのポインタを持って観測する。
type SelectStageCanvas struct {
	scene    *scene.SelectStageScene
	canvas   *canvas.Canvas
	routeMap *ui.RouteMapUI[stageinfo.Location]
	guide    *ui.Text

	// labels は駅の表示名。案内文の組み立てに使う。
	labels map[stageinfo.Location]string
	// shownRevision は反映済みの SelectStageScene.SelectionRevision。
	shownRevision int
	// shownCurrent は反映済みの現在地。
	shownCurrent stageinfo.Location
}

// NewSelectStageCanvas は SelectStageCanvas を生成する。
func NewSelectStageCanvas(s *scene.SelectStageScene, gdata *game.Gamedata) (*SelectStageCanvas, error) {
	if s == nil {
		return nil, fmt.Errorf("SelectStageScene が nil です")
	}
	if gdata == nil {
		return nil, fmt.Errorf("gdata が nil です")
	}
	if gdata.Font == nil {
		return nil, fmt.Errorf("フォントが nil です")
	}

	m := s.StageMap()

	// 描画順を安定させるため、駅は Location の昇順に並べる
	nodes := []ui.RouteMapNode[stageinfo.Location]{}
	labels := make(map[stageinfo.Location]string, len(m.Stations))
	for _, loc := range slices.Sorted(maps.Keys(m.Stations)) {
		station := m.Stations[loc]
		labels[loc] = station.Label
		nodes = append(nodes, ui.RouteMapNode[stageinfo.Location]{
			ID:    loc,
			Pos:   station.Pos,
			Label: station.Label,
		})
	}
	edges := make([][2]stageinfo.Location, 0, len(m.Lines))
	for _, line := range m.Lines {
		edges = append(edges, [2]stageinfo.Location{line.From, line.To})
	}

	routeMap, err := ui.NewRouteMapUI(utils.Vec2{X: m.Width, Y: m.Height}, nodes, edges, gdata.Font)
	if err != nil {
		return nil, fmt.Errorf("ui.NewRouteMapUI() でエラーが発生しました。\n%w", err)
	}

	guide, err := ui.NewText(utils.Position{}, color.Orange(), "", ui.FontFace(gdata.Font, 12))
	if err != nil {
		return nil, fmt.Errorf("ui.NewText() でエラーが発生しました。\n%w", err)
	}

	c := canvas.NewCanvas()
	for _, item := range []canvas.CanvasItem{routeMap, guide} {
		if err := c.Add(item); err != nil {
			return nil, fmt.Errorf("ステージ選択画面のCanvas構築に失敗しました。\n%w", err)
		}
	}

	sc := &SelectStageCanvas{scene: s, canvas: c, routeMap: routeMap, guide: guide}
	sc.syncWithScene()
	return sc, nil
}

// syncWithScene は Scene の選択中の駅・現在地を路線図と案内テキストへ反映する。
func (c *SelectStageCanvas) syncWithScene() {
	selected := c.scene.SelectedLocation()
	c.routeMap.SetSelected(selected)
	c.routeMap.SetMarker(c.scene.CurrentLocation())
	c.guide.SetText("行き先を選択\nWASD / 矢印キー: 移動　Enter: 決定\n行き先: " + c.scene.StageMap().Stations[selected].Label)
}

// Update implements [scene.SelectStageCanvasAPI].
func (c *SelectStageCanvas) Update(active bool) error {
	c.syncWithScene()
	return c.canvas.Update(active)
}

// Draw implements [scene.SelectStageCanvasAPI].
func (c *SelectStageCanvas) Draw(screen *ebiten.Image) {
	c.canvas.Draw(screen)
}

// IsSelectStageCanvas implements [scene.SelectStageCanvasAPI].
func (c *SelectStageCanvas) IsSelectStageCanvas() {}

var _ scene.SelectStageCanvasAPI = (*SelectStageCanvas)(nil)
