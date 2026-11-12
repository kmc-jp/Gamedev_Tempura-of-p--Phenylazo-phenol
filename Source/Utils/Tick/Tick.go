package tick

type Timer struct {
	//unit: 0.1 second
	//0以外の場合は計測中として扱う
	setTime uint16

	tick      uint64
	startTick uint64
	expired   bool

	suspend bool
}

var timers = map[*Timer]struct{}{}

// Tick advances every timer created by NewTimer by one game frame.
func Tick() {
	for timer := range timers {
		if timer.suspend {
			continue
		}
		timer.tick++
		if timer.setTime != 0 && timer.tick >= timer.startTick+uint64(timer.setTime)*6 {
			timer.expired = true
			timer.setTime = 0
		}
	}
}

func NewTimer() *Timer {
	timer := &Timer{}
	timers[timer] = struct{}{}
	return timer
}

// Expired reports whether the timer reached its configured duration.
// The notification is consumed when it returns true.
func (timer *Timer) Expired() bool {
	if !timer.expired {
		return false
	}
	timer.expired = false
	return true
}

// Set configures the timer duration in 0.1-second units.
func (timer *Timer) Set(time uint16) error {
	if time == 0 {
		timer.setTime = 0
		timer.expired = false
		return nil
	}

	//すでに計測中の場合はTimer設定をスキップ
	if timer.setTime != 0 {
		return nil
	}
	timer.setTime = time
	timer.startTick = timer.tick
	return nil
}

func (timer *Timer) Suspend() {
	timer.suspend = true
}

func (timer *Timer) Resume() {
	timer.suspend = false
}
