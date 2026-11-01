package tick

type Timer struct {
	//unit: 0.1 second
	//0以外の場合は計測中として扱う
	Time uint16

	tick      uint64
	startTick uint64
	expired   bool
}

var timers = map[*Timer]struct{}{}

// Tick advances every timer created by NewTimer by one game frame.
func Tick() {
	for timer := range timers {
		timer.tick++
		if timer.Time != 0 && timer.tick >= timer.startTick+uint64(timer.Time)*6 {
			timer.expired = true
			timer.Time = 0
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
		timer.Time = 0
		timer.expired = false
		return nil
	}

	//すでに計測中の場合はTimer設定をスキップ
	if timer.Time != 0 {
		return nil
	}
	timer.Time = time
	timer.startTick = timer.tick
	return nil
}
