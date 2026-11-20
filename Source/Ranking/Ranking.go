// Package ranking はスコアランキングの保存と読み込みを行う。ebitenには依存しない。
package ranking

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	// MaxEntries はランキングに残す件数。
	MaxEntries = 10
	// AppearanceSize は見た目enumの個数。
	AppearanceSize = 7
	// FileName は保存ファイル名。
	FileName = "ranking.json"
	version  = 1
)

// PartRecord は組み立てに使ったパーツ1つ分の記録。
type PartRecord struct {
	Kind uint   `json:"kind"`
	Name string `json:"name"`
	Rank uint   `json:"rank"`
}

// Entry はランキング1件分。
type Entry struct {
	Score    int          `json:"score"`
	PlayedAt time.Time    `json:"playedAt"`
	Parts    []PartRecord `json:"parts"`
	// Appearance は見た目のenum。順序は caseBack, mb, cooler, gpu, caseLid, gaming, sticker。
	Appearance [AppearanceSize]uint `json:"appearance"`
}

// Ranking はスコア降順のランキング。
type Ranking struct {
	Version int     `json:"version"`
	Entries []Entry `json:"entries"`
}

// New は空のランキングを返す。
func New() *Ranking {
	return &Ranking{Version: version}
}

// Insert は e を挿入し、入った順位(0始まり)を返す。圏外なら -1。
// 同点の場合は既存の記録を上位とする。
func (r *Ranking) Insert(e Entry) int {
	pos := len(r.Entries)
	for i, existing := range r.Entries {
		if e.Score > existing.Score {
			pos = i
			break
		}
	}
	if pos >= MaxEntries {
		return -1
	}
	r.Entries = append(r.Entries, Entry{})
	copy(r.Entries[pos+1:], r.Entries[pos:])
	r.Entries[pos] = e
	if len(r.Entries) > MaxEntries {
		r.Entries = r.Entries[:MaxEntries]
	}
	return pos
}

// Load は path からランキングを読む。ファイルが無ければ空のランキングを返す。
func Load(path string) (*Ranking, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return New(), nil
		}
		return New(), fmt.Errorf("ランキングファイルの読み込みに失敗しました: %w", err)
	}
	r := New()
	if err := json.Unmarshal(data, r); err != nil {
		return New(), fmt.Errorf("ランキングファイルの解析に失敗しました: %w", err)
	}
	if len(r.Entries) > MaxEntries {
		r.Entries = r.Entries[:MaxEntries]
	}
	return r, nil
}

// Save は一時ファイルに書いてから置き換え、書き込み途中の破損を避ける。
func (r *Ranking) Save(path string) error {
	r.Version = version
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return fmt.Errorf("ランキングのJSON化に失敗しました: %w", err)
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return fmt.Errorf("ランキングの書き込みに失敗しました: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("ランキングファイルの置き換えに失敗しました: %w", err)
	}
	return nil
}

// DefaultPath は実行ファイルと同じフォルダの ranking.json を返す。
// go run のように実行ファイルが一時フォルダにある場合はカレントディレクトリを使う。
func DefaultPath() (string, error) {
	exe, err := os.Executable()
	if err == nil {
		dir := filepath.Dir(exe)
		if !isUnder(dir, os.TempDir()) {
			return filepath.Join(dir, FileName), nil
		}
	}
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("保存先フォルダの取得に失敗しました: %w", err)
	}
	return filepath.Join(wd, FileName), nil
}

func isUnder(path, base string) bool {
	rel, err := filepath.Rel(base, path)
	if err != nil {
		return false
	}
	return rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

// Store はファイルに紐づいたランキングの読み書き窓口。
type Store struct {
	path string
}

// NewStore は保存先 path のStoreを返す。
func NewStore(path string) *Store {
	return &Store{path: path}
}

// NewDefaultStore は DefaultPath を保存先にしたStoreを返す。
func NewDefaultStore() (*Store, error) {
	path, err := DefaultPath()
	if err != nil {
		return nil, err
	}
	return NewStore(path), nil
}

// Entries は現在の記録を返す。読み込みに失敗した場合は空と err を返す。
func (s *Store) Entries() ([]Entry, error) {
	r, err := Load(s.path)
	return r.Entries, err
}

// Record は e を記録して保存し、順位(0始まり、圏外は -1)を返す。
// 既存ファイルが壊れていた場合は空から作り直す。保存に失敗しても順位は返す。
func (s *Store) Record(e Entry) (int, error) {
	r, loadErr := Load(s.path)
	rank := r.Insert(e)
	if err := r.Save(s.path); err != nil {
		return rank, err
	}
	return rank, loadErr
}
