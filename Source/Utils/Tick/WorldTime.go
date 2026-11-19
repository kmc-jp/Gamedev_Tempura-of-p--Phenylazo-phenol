package tick

import "fmt"

type WorldTimer struct {
	Timer

	//ゲームスタートからの経過時間
	passedTime uint64 //unit: min

	isPM bool
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
	timer := &WorldTimer{}
	timers[&timer.Timer] = struct{}{}
	timer.Set(PlayableTime * CountSpeed)
	return timer
}

func (wtimer *WorldTimer) syncTime() {
	wtimer.passedTime = wtimer.GetCurrentTick() / countSpeed_tick
}

func (wtimer *WorldTimer) setIsPM() {
	departureMin := Departuretime * 60
	if uint64(departureMin)+wtimer.passedTime >= 12*60 {
		wtimer.isPM = true
	}
}

func (wtimer *WorldTimer) GetPassedTimeInMinite() uint64 {
	wtimer.syncTime()
	return wtimer.passedTime
}

// GetTimeText はゲーム内時刻を hh:mm 形式で返す。
func (wtimer *WorldTimer) GetTimeText() string {
	totalMinutes := uint64(Departuretime*60) + wtimer.GetPassedTimeInMinite()
	endMinutes := uint64(Endtime * 60)
	if totalMinutes > endMinutes {
		totalMinutes = endMinutes
	}
	if !wtimer.isPM {
		wtimer.setIsPM()
	}
	ampm := ""
	switch wtimer.isPM {
	case true:
		ampm = "PM"
	case false:
		ampm = "AM"
	}
	return fmt.Sprintf("%s  %02d:%02d", ampm, totalMinutes/60, totalMinutes%60)
}
