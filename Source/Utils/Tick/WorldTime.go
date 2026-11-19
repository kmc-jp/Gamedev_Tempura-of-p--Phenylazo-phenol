package tick

import "fmt"

type WorldTimer struct {
	Timer

	//ゲームスタートからの経過時間
	passedTime uint64 //unit: min
}

// CountSpeed 経過ごとに WorldTime 1 min 経過
const (
	//画面上の時計表示
	Departuretime = 10
	Endtime       = 16

	playableTime_hour = Endtime - Departuretime
	PlayableTime      = 60 * playableTime_hour //分単位のプレイ可能時間

	CountSpeed      = 1 //WorldTime 1 min あたりの経過時間 0.1sec
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

func (wtimer *WorldTimer) GetPassedTimeInMinite() uint64 {
	wtimer.syncTime()
	return wtimer.passedTime
}

func formatWorldTime(totalMinutes uint64) string {
	endMinutes := uint64(Endtime * 60)
	if totalMinutes > endMinutes {
		totalMinutes = endMinutes
	}

	hour24 := totalMinutes / 60
	minute := totalMinutes % 60
	period := "午前"
	hour12 := hour24
	if hour24 >= 12 {
		period = "午後"
		hour12 = hour24 - 12
	}
	if hour12 == 0 {
		hour12 = 12
	}
	return fmt.Sprintf("%s %02d:%02d", period, hour12, minute)
}

// GetTimeText はゲーム内時刻を日本式の午前・午後付き 12時間表記で返す。
func (wtimer *WorldTimer) GetTimeText() string {
	totalMinutes := uint64(Departuretime*60) + wtimer.GetPassedTimeInMinite()
	return formatWorldTime(totalMinutes)
}
