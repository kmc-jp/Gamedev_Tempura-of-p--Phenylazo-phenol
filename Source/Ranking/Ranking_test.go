package ranking

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInsertOrderAndLimit(t *testing.T) {
	r := New()
	for i := 1; i <= 12; i++ {
		r.Insert(Entry{Score: i * 100})
	}
	if len(r.Entries) != MaxEntries {
		t.Fatalf("件数 = %d, want %d", len(r.Entries), MaxEntries)
	}
	if r.Entries[0].Score != 1200 || r.Entries[MaxEntries-1].Score != 300 {
		t.Fatalf("並びが不正: %v ... %v", r.Entries[0].Score, r.Entries[MaxEntries-1].Score)
	}
	if got := r.Insert(Entry{Score: 100}); got != -1 {
		t.Fatalf("圏外の順位 = %d, want -1", got)
	}
}

func TestInsertTieKeepsExistingAbove(t *testing.T) {
	r := New()
	r.Insert(Entry{Score: 500, Parts: []PartRecord{{Name: "old"}}})
	got := r.Insert(Entry{Score: 500, Parts: []PartRecord{{Name: "new"}}})
	if got != 1 || r.Entries[0].Parts[0].Name != "old" {
		t.Fatalf("同点は既存が上位のはず: rank=%d", got)
	}
}

func TestInsertTieWhenFullIsOutOfRange(t *testing.T) {
	r := New()
	for range MaxEntries {
		r.Insert(Entry{Score: 500})
	}
	if got := r.Insert(Entry{Score: 500}); got != -1 {
		t.Fatalf("満杯で同点の順位 = %d, want -1", got)
	}
}

func TestSaveLoadRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), FileName)
	r := New()
	r.Insert(Entry{
		Score:      1234,
		Parts:      []PartRecord{{Kind: 1, Name: "GPU", Rank: 3}},
		Appearance: [AppearanceSize]uint{1, 2, 3, 4, 5, 6, 7},
	})
	if err := r.Save(path); err != nil {
		t.Fatal(err)
	}
	loaded, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded.Entries) != 1 || loaded.Entries[0].Score != 1234 ||
		loaded.Entries[0].Appearance != r.Entries[0].Appearance ||
		loaded.Entries[0].Parts[0].Name != "GPU" {
		t.Fatalf("往復で内容が変わった: %+v", loaded.Entries)
	}
}

func TestLoadMissingFileIsEmpty(t *testing.T) {
	r, err := Load(filepath.Join(t.TempDir(), "none.json"))
	if err != nil || len(r.Entries) != 0 {
		t.Fatalf("err=%v entries=%d", err, len(r.Entries))
	}
}

func TestStoreRecoversFromCorruptFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), FileName)
	if err := os.WriteFile(path, []byte("{broken"), 0o644); err != nil {
		t.Fatal(err)
	}
	s := NewStore(path)
	rank, err := s.Record(Entry{Score: 10})
	if rank != 0 {
		t.Fatalf("rank = %d, want 0", rank)
	}
	if err == nil {
		t.Fatal("壊れたファイルは警告用にエラーを返すはず")
	}
	entries, err := s.Entries()
	if err != nil || len(entries) != 1 {
		t.Fatalf("作り直されていない: err=%v n=%d", err, len(entries))
	}
}
