package tick

import "fmt"

//Timerを使う場合はUpdate()にCount()処理が必要
type Timer struct {
	//unit: 0.1 second
	//0以外の場合は計測中として扱う
	Time uint16

	tick      uint16
	startTick uint16
}

const maxTick uint16 = 36000 // <= up to 65535

// Timerを使う場合はUpdate()にこの処理が必要
func (timer *Timer) Count() bool {
	counted := false
	timer.tick++
	if timer.Time != 0 && timer.Time <= maxTick/6 {
		ct := timer.startTick + timer.Time*6
		if ct >= maxTick {
			ct -= maxTick
		}
		if timer.tick >= ct {
			counted = true
			timer.Time = 0
		}
	}
	if timer.tick >= maxTick {
		timer.tick = 0
	}
	return counted
}

func NewTimer() *Timer {
	return &Timer{}
}

//SetTimer(カウントする時間(0.1 sec))
func (timer *Timer) SetTimer(time uint16) error {
	//すでに計測中の場合はTimer設定をスキップ
	if timer.Time != 0 {
		return nil
	}
	if time > maxTick/6 {
		return fmt.Errorf("SetTimer() での計測可能時間を超過しています")
	}
	timer.Time = time
	timer.startTick = timer.tick
	return nil
}
