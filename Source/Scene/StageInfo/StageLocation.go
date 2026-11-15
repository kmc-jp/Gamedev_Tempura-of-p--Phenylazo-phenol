package stageinfo

// Location は路線図上の駅を表す。
type Location uint

const (
	Red Location = iota
	Blue
	Green
	Orange
	Purple

	// Departure は自宅駅（ゲーム開始時の現在地）。
	Departure
)

// Location2String は Location を TOML などで使う駅名へ変換する。
var Location2String = map[Location]string{
	Red:       "Red",
	Blue:      "Blue",
	Green:     "Green",
	Orange:    "Orange",
	Purple:    "Purple",
	Departure: "Departure",
}

// String2Location は駅名から Location を引く。
var String2Location = map[string]Location{
	"Red":       Red,
	"Blue":      Blue,
	"Green":     Green,
	"Orange":    Orange,
	"Purple":    Purple,
	"Departure": Departure,
}
