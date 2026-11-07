package Utils

//only linear motion
const (
	Move_Walk_speedBased motiontype = iota
	Move_Walk_timeBased
	Move_Decelerate

	Move_Finished
	Moving
)

type motiontype uint

type PathAnimateManager struct {
	currentStep uint
	ToStep      uint

	direction    Vec2
	destinaton   Vec2
	speed        float64
	acceleration float64
}

func NewPathAnimateManager() *PathAnimateManager {
	return &PathAnimateManager{
		currentStep: 0,
		ToStep:      0,

		direction:    NewZeroVec2(),
		destinaton:   NewZeroVec2(),
		speed:        0,
		acceleration: 0,
	}
}

func (pam *PathAnimateManager) SetPathAnimation(atype motiontype, to Vec2, fromSpeed float64, time uint, p *Position) error {
	if pam.currentStep != 0 {
		return nil
	}
	currentPos := p.PtoV()
	distance := currentPos.Distance(to)
	direction, err := UnitDirection(p.PtoV(), to)
	if err != nil {
		//零ベクトルの正規化エラー
		return err
	}

	switch atype {
	case Move_Walk_speedBased:
		distParStep := fromSpeed / 60
		StepsRequired := uint(distance / distParStep)

		pam.currentStep = 0
		pam.ToStep = StepsRequired
		pam.direction = direction
		pam.destinaton = to
		pam.speed = fromSpeed
		pam.acceleration = 0

	case Move_Walk_timeBased:
		StepsRequired := time * 6
		distParStep := distance / float64(StepsRequired)
		spd := distParStep * 60

		pam.currentStep = 0
		pam.ToStep = StepsRequired
		pam.direction = direction
		pam.destinaton = to
		pam.speed = spd
		pam.acceleration = 0
	case Move_Decelerate:
		StepsRequired := time * 6
		acl := (distance - fromSpeed*float64(StepsRequired)) / (float64(StepsRequired) * float64(StepsRequired))

		pam.currentStep = 0
		pam.ToStep = StepsRequired
		pam.direction = direction
		pam.destinaton = to
		pam.speed = fromSpeed
		pam.acceleration = acl
	}
	return nil
}

func (p *Position) PathAnimate(pam *PathAnimateManager) motiontype {
	// position += direction norm * (v/60 * t + a/60 * t * t)
	p.Move(pam.direction.ScalarProd(float64(pam.speed*float64(pam.currentStep) + pam.acceleration*float64(pam.currentStep)*float64(pam.currentStep))).ScalarProd(1.0 / 60.0))
	pam.currentStep++
	if pam.currentStep == pam.ToStep {
		p.SetPosition(pam.destinaton)
		return Move_Finished
	}
	return Moving
}
