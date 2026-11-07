package Utils

import (
	"fmt"
	"math"
)

const (
	ticksPerSecond     = 60
	ticksPerDecisecond = ticksPerSecond / 10
)

// only linear motion
const (
	Move_Walk_speedBased Motiontype = iota
	Move_Walk_timeBased
	Move_Decelerate

	Move_Finished
	Moving
)

type Motiontype uint

type PathAnimateManager struct {
	currentStep uint
	ToStep      uint

	active       bool
	start        Vec2
	direction    Vec2
	destinaton   Vec2
	speed        float64
	acceleration float64
}

func NewPathAnimateManager() *PathAnimateManager {
	return &PathAnimateManager{
		currentStep: 0,
		ToStep:      0,

		active:       false,
		start:        NewZeroVec2(),
		direction:    NewZeroVec2(),
		destinaton:   NewZeroVec2(),
		speed:        0,
		acceleration: 0,
	}
}

func (pam *PathAnimateManager) SetPathAnimation(atype Motiontype, to Vec2, fromSpeed float64, time uint, p *Position) error {
	if p == nil {
		return fmt.Errorf("position が nil です")
	}
	if pam.active {
		return fmt.Errorf("アニメーションは実行中です")
	}
	currentPos := p.PtoV()
	distance := currentPos.Distance(to)
	if distance == 0 {
		pam.reset()
		return nil
	}

	switch atype {
	case Move_Walk_speedBased:
		if fromSpeed <= 0 {
			return fmt.Errorf("速度は 0 より大きい必要があります: %v", fromSpeed)
		}
		distParStep := fromSpeed / ticksPerSecond
		steps := uint(math.Ceil(distance / distParStep))
		if steps == 0 {
			steps = 1
		}

		pam.configure(currentPos, to, steps)
		pam.speed = fromSpeed
		pam.acceleration = 0

	case Move_Walk_timeBased:
		steps, err := animationSteps(time)
		if err != nil {
			return err
		}

		pam.configure(currentPos, to, steps)
		pam.speed = distance / float64(steps) * ticksPerSecond
		pam.acceleration = 0
	case Move_Decelerate:
		if fromSpeed < 0 {
			return fmt.Errorf("初速度は 0 以上である必要があります: %v", fromSpeed)
		}
		steps, err := animationSteps(time)
		if err != nil {
			return err
		}

		pam.configure(currentPos, to, steps)
		pam.speed = fromSpeed
		// s = v0*t + 1/2*a*t^2. The time unit of one step is 1/60 s.
		totalTime := float64(steps) / ticksPerSecond
		pam.acceleration = 2 * (distance - fromSpeed*totalTime) / (totalTime * totalTime)
		finalSpeed := fromSpeed + pam.acceleration*totalTime
		if pam.acceleration >= 0 || finalSpeed < 0 {
			pam.reset()
			return fmt.Errorf("指定された距離、初速度、時間では減速運動を成立させられません")
		}
	default:
		return fmt.Errorf("未対応のモーション種別です: %d", atype)
	}
	return nil
}

func (p *Position) PathAnimate(pam *PathAnimateManager) Motiontype {
	if pam == nil || !pam.active {
		return Move_Finished
	}

	// Use the position at the next tick so the first call moves immediately.
	pam.currentStep++
	t := float64(pam.currentStep) / ticksPerSecond
	distance := pam.speed*t + 0.5*pam.acceleration*t*t
	p.SetPosition(pam.start.Add(pam.direction.ScalarProd(distance)))
	if pam.currentStep >= pam.ToStep {
		p.SetPosition(pam.destinaton)
		pam.reset()
		return Move_Finished
	}
	return Moving
}

func animationSteps(time uint) (uint, error) {
	if time == 0 {
		return 0, fmt.Errorf("アニメーション時間は 0 より大きい必要があります")
	}
	return time * ticksPerDecisecond, nil
}

func (pam *PathAnimateManager) configure(start, destination Vec2, steps uint) {
	direction, _ := UnitDirection(start, destination)
	pam.currentStep = 0
	pam.ToStep = steps
	pam.active = true
	pam.start = start
	pam.direction = direction
	pam.destinaton = destination
}

func (pam *PathAnimateManager) reset() {
	pam.currentStep = 0
	pam.ToStep = 0
	pam.active = false
	pam.start = NewZeroVec2()
	pam.direction = NewZeroVec2()
	pam.destinaton = NewZeroVec2()
	pam.speed = 0
	pam.acceleration = 0
}
