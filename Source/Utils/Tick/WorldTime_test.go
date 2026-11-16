package tick

import "testing"

func TestWorldTimerAdvancesAndCanBeSuspended(t *testing.T) {
	timer := NewWorldTimer()

	if got := timer.GetTimeText(); got != "10:00" {
		t.Fatalf("初期時刻 = %q, want %q", got, "10:00")
	}

	for range countSpeed_tick {
		Tick()
	}
	if got := timer.GetTimeText(); got != "10:01" {
		t.Fatalf("1分経過後の時刻 = %q, want %q", got, "10:01")
	}

	timer.Suspend()
	for range countSpeed_tick {
		Tick()
	}
	if got := timer.GetTimeText(); got != "10:01" {
		t.Fatalf("停止中の時刻 = %q, want %q", got, "10:01")
	}

	timer.Resume()
	for range countSpeed_tick {
		Tick()
	}
	if got := timer.GetTimeText(); got != "10:02" {
		t.Fatalf("再開後の時刻 = %q, want %q", got, "10:02")
	}
}
