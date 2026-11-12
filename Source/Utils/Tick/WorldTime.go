package tick

type WorldTimer struct {
	Timer

	passedTime uint64 //unit: min
}

// CountSpeed 経過ごとに WorldTime 1 min 経過
const (
	//画面上の時計表示
	Departuretime = 10
	Endtime       = 16

	playableTime_hour = Endtime - Departuretime
	PlayableTime      = 60 * playableTime_hour //分単位のプレイ可能時間

	CountSpeed      = 20 //WorldTime 1 min あたりの経過時間 0.1sec
	countSpeed_tick = CountSpeed * 6
)

func NewWorldTimer() *WorldTimer {
	return &WorldTimer{Timer: *NewTimer()}
}

func (wtimer *WorldTimer) syncTime() {
	wtimer.passedTime = wtimer.GetCurrentTick() / countSpeed_tick
}

func (wtimer *WorldTimer) GetPassedTimeInMinite() uint64 {
	wtimer.syncTime()
	return wtimer.passedTime
}
