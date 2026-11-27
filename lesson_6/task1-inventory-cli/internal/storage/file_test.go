package storage

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Anton-Babaev/vk_homework/lesson_6/task1-inventory-cli/internal/inventory"
	"github.com/Anton-Babaev/vk_homework/lesson_6/task1-inventory-cli/internal/item"
)

// TestFileStorage_RoundTripNoTempLeftover проверяет сохранение и атомарную
// запись: после Save не должно оставаться временных файлов .inventory-*.
func TestFileStorage_RoundTripNoTempLeftover(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "inv.txt")

	s := NewFileStorage(path)
	inv := inventory.New()
	inv.AddItem(item.NewWeapon("Меч", 2.5, 10, 3))

	if err := s.Save(inv); err != nil {
		t.Fatalf("Save: %v", err)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("в каталоге %d файлов, ожидался 1 (без временных)", len(entries))
	}

	got, err := s.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.Len() != 1 {
		t.Fatalf("предметов %d, want 1", got.Len())
	}
}

func TestFileStorage_MissingFileIsEmpty(t *testing.T) {
	s := NewFileStorage(filepath.Join(t.TempDir(), "nope.txt"))
	inv, err := s.Load()
	if err != nil {
		t.Fatalf("отсутствующий файл не должен быть ошибкой: %v", err)
	}
	if inv.Len() != 0 {
		t.Fatalf("предметов %d, want 0", inv.Len())
	}
}

func TestFileStorage_BrokenLineError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "inv.txt")
	if err := os.WriteFile(path, []byte("Weapon|1|Меч|10|3|2.5\ngarbage\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	s := NewFileStorage(path)
	inv, err := s.Load()
	if err == nil {
		t.Fatal("ожидалась ошибка о битой строке")
	}
	var pie *inventory.ParseItemError
	if !errors.As(err, &pie) {
		t.Fatalf("ожидалась ParseItemError, получено %T", err)
	}
	if !strings.Contains(err.Error(), "битая строка") {
		t.Fatalf("ошибка должна быть обёрнута как «битая строка»: %v", err)
	}
	if inv.Len() != 1 {
		t.Fatalf("валидных предметов %d, want 1", inv.Len())
	}
}
