package ui

import (
	canvas "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Canvas"
	uiutils "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/UI/UIUtils"
	utils "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils"
	color "Gamedev_Tempura-of-p--Phenylazo-phenol/Source/Utils/Color"
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// 基準サイズ上の大きさ。描画時に画面へのフィット倍率が掛かる。
const (
	routeNodeRadius   = 8.0
	routeMarkerRadius = 4.0
	routeLineWidth    = 3.0
	routeLabelSize    = 10.0
	routeRingWidth    = 2.0
	routeRingMargin   = 3.0
)

// RouteMapNode は路線図上のノード（駅）である。
type RouteMapNode[N comparable] struct {
	ID    N
	Pos   utils.Position // 基準サイズ上の座標（ノードの中心）
	Label string
	Image *ebiten.Image // nil なら円で描く。画像があれば中心を Pos に合わせて描く
}

// RouteMapUI は、ノードとそれを結ぶ辺を路線図として描画する。
// 選択中のノードと現在地のノードを強調表示する。入力は扱わず、
// どれを選択中にするかは SetSelected で外から渡す。
//
// 座標は refSize を基準とし、画面サイズに合わせてアスペクト比を保ったまま拡大縮小して中央に描く。
type RouteMapUI[N comparable] struct {
	refSize    utils.Vec2
	nodes      []RouteMapNode[N]
	index      map[N]int
	edges      [][2]N
	fontSource *text.GoTextFaceSource
	background *ebiten.Image // 任意。基準サイズ全体に引き伸ばして描く

	selected    N
	hasSelected bool
	marker      N
	hasMarker   bool

	frame int // 選択リングの脈動用。Update で進める
}

// NewRouteMapUI は RouteMapUI を生成する。
func NewRouteMapUI[N comparable](refSize utils.Vec2, nodes []RouteMapNode[N], edges [][2]N, fontSource *text.GoTextFaceSource) (*RouteMapUI[N], error) {
	if refSize.X <= 0 || refSize.Y <= 0 {
		return nil, fmt.Errorf("refSize は正である必要があります: %v", refSize)
	}
	if fontSource == nil {
		return nil, fmt.Errorf("fontSource が nil です")
	}

	index := make(map[N]int, len(nodes))
	for i, node := range nodes {
		if _, dup := index[node.ID]; dup {
			return nil, fmt.Errorf("ノードID %v が重複しています", node.ID)
		}
		index[node.ID] = i
	}
	for _, edge := range edges {
		for _, end := range edge {
			if _, ok := index[end]; !ok {
				return nil, fmt.Errorf("辺の端 %v に対応するノードがありません", end)
			}
		}
	}

	return &RouteMapUI[N]{
		refSize:    refSize,
		nodes:      nodes,
		index:      index,
		edges:      edges,
		fontSource: fontSource,
	}, nil
}

// SetSelected は選択中のノードを設定する。未知のIDなら何もしない。
func (u *RouteMapUI[N]) SetSelected(id N) {
	if _, ok := u.index[id]; !ok {
		return
	}
	u.selected = id
	u.hasSelected = true
}

// SetMarker は現在地のノードを設定する。未知のIDなら何もしない。
func (u *RouteMapUI[N]) SetMarker(id N) {
	if _, ok := u.index[id]; !ok {
		return
	}
	u.marker = id
	u.hasMarker = true
}

// SetBackground は路線図の背景画像を設定する。nil で背景なしに戻す。
func (u *RouteMapUI[N]) SetBackground(img *ebiten.Image) {
	u.background = img
}

// Update は選択リングの脈動のためにフレームを進める。
func (u *RouteMapUI[N]) Update(active bool) error {
	if active {
		u.frame++
	}
	return nil
}

// Draw は路線図を描画する。
func (u *RouteMapUI[N]) Draw(screen *ebiten.Image) {
	b := screen.Bounds()
	scale := math.Min(float64(b.Dx())/u.refSize.X, float64(b.Dy())/u.refSize.Y)
	offX := float64(b.Min.X) + (float64(b.Dx())-u.refSize.X*scale)/2
	offY := float64(b.Min.Y) + (float64(b.Dy())-u.refSize.Y*scale)/2
	toScreen := func(p utils.Position) (float32, float32) {
		return float32(offX + p.X*scale), float32(offY + p.Y*scale)
	}

	if u.background != nil {
		bg := u.background.Bounds()
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(u.refSize.X*scale/float64(bg.Dx()), u.refSize.Y*scale/float64(bg.Dy()))
		op.GeoM.Translate(offX, offY)
		screen.DrawImage(u.background, op)
	}

	/*
		lineColor := color.LightGray().ToRGBA()
		for _, edge := range u.edges {
			x0, y0 := toScreen(u.nodes[u.index[edge[0]]].Pos)
			x1, y1 := toScreen(u.nodes[u.index[edge[1]]].Pos)
			vector.StrokeLine(screen, x0, y0, x1, y1, float32(routeLineWidth*scale), lineColor, true)
		}
	*/

	nodeColor := color.White().ToRGBA()
	markerColor := color.RoyalBlue().ToRGBA()
	for _, node := range u.nodes {
		cx, cy := toScreen(node.Pos)
		if node.Image != nil {
			size := node.Image.Bounds().Size()
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(-float64(size.X)/2, -float64(size.Y)/2)
			op.GeoM.Scale(scale, scale)
			op.GeoM.Translate(float64(cx), float64(cy))
			screen.DrawImage(node.Image, op)
		} else {
			vector.FillCircle(screen, cx, cy, float32(routeNodeRadius*scale), nodeColor, true)
		}

		if u.hasMarker && node.ID == u.marker {
			vector.FillCircle(screen, cx, cy, float32(routeMarkerRadius*scale), markerColor, true)
		}

		if u.hasSelected && node.ID == u.selected {
			r := (routeNodeRadius + routeRingMargin + math.Sin(float64(u.frame)/8)) * scale
			vector.StrokeCircle(screen, cx, cy, float32(r), float32(routeRingWidth*scale), uiutils.SelectionBorderColor, true)
		}

		u.drawLabel(screen, node.Label, cx, cy, scale)
	}
}

// drawLabel はノードの下に中央揃えでラベルを描く。
func (u *RouteMapUI[N]) drawLabel(screen *ebiten.Image, label string, cx, cy float32, scale float64) {
	if label == "" {
		return
	}
	face := FontFace(u.fontSource, routeLabelSize*scale)
	w, _ := text.Measure(label, face, 0)
	pos := utils.Position{
		X: float64(cx) - w/2,
		Y: float64(cy) + (routeNodeRadius+routeRingMargin+2)*scale,
	}
	DrawText(pos, color.White(), label, face, screen)
}

var _ canvas.CanvasItem = (*RouteMapUI[int])(nil)
