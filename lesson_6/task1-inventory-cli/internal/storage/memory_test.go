package storage

import (
	"errors"
	"strings"
	"testing"

	"github.com/Anton-Babaev/vk_homework/lesson_6/task1-inventory-cli/internal/inventory"
	"github.com/Anton-Babaev/vk_homework/lesson_6/task1-inventory-cli/internal/item"
)

func TestMemoryStorage_RoundTrip(t *testing.T) {
	s := NewMemoryStorage()
	inv := inventory.New()
	inv.AddItem(item.NewWeapon("Меч", 2.5, 10, 3))

	if err := s.Save(inv); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := s.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.Len() != 1 {
		t.Fatalf("предметов %d, want 1", got.Len())
	}
}

// TestMemoryStorage_BrokenLineError проверяет, что in-memory хранилище
// обрабатывает битую строку так же, как файловое (бонус ДЗ6).
func TestMemoryStorage_BrokenLineError(t *testing.T) {
	s := NewMemoryStorage()
	// Восстанавливаем ситуацию «битая запись в сохранённых данных».
	s.data = []byte("Weapon|1|Меч|10|3|2.5\ngarbage\n")

	_, err := s.Load()
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
}
