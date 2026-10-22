package tick

//Timerを使う場合はUpdate()にCount()処理が必要
type Timer struct {
	//unit: 0.1 second
	//0以外の場合は計測中として扱う
	Time uint

	tick      uint
	startTick uint
}

// Timerを使う場合はUpdate()にこの処理が必要
func (timer *Timer) Count() bool {
	timer.tick++
	if timer.startTick+timer.Time*6 <= timer.tick {
		timer.startTick = timer.tick
		return true
	}

	return false
}

func NewTimer() *Timer {
	return &Timer{}
}

//SetTimer(カウントする時間(0.1 sec) 停止する場合は SetTimer(0) を実行
func (timer *Timer) Set(time uint) {
	//すでに計測中の場合はTimer設定をスキップ
	if timer.Time != 0 {
		return
	}
	timer.Time = time
	timer.startTick = timer.tick
}

func (timer *Timer) Stop() {
	timer.Time = 0
}
