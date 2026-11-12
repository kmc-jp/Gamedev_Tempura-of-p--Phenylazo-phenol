package game

var WorldTime uint

// CountSpeed 経過ごとに WorldTime 1 min 経過
const (
	//画面上の時計表示
	Departuretime = 10
	Endtime       = 16

	PlayableTime_hour = Endtime - Departuretime
	PlayableTime      = 60 * PlayableTime_hour //分単位のプレイ可能時間

	CountSpeed = 20 //WorldTime 1 min あたりの経過時間 0.1sec
)
