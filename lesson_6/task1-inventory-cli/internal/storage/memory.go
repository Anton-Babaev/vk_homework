package storage

import (
	"bytes"

	"github.com/Anton-Babaev/vk_homework/lesson_6/task1-inventory-cli/internal/inventory"
)

// MemoryStorage хранит инвентарь в памяти в том же текстовом формате,
// что и FileStorage. Это гарантирует одинаковые правила загрузки,
// сохранения и обработки ошибок: обе реализации используют
// inventory.Load/Save, отличаются только носителем данных.
type MemoryStorage struct {
	data []byte
}

// NewMemoryStorage создаёт пустое in-memory хранилище.
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{}
}

// Load читает инвентарь из сохранённого в памяти снимка.
func (s *MemoryStorage) Load() (*inventory.Inventory, error) {
	inv := inventory.New()
	if len(s.data) == 0 {
		return inv, nil
	}
	if err := inv.Load(bytes.NewReader(s.data)); err != nil {
		return inv, wrapLoadError("", err)
	}
	return inv, nil
}

// Save сериализует инвентарь в память.
func (s *MemoryStorage) Save(inv *inventory.Inventory) error {
	var buf bytes.Buffer
	if err := inv.Save(&buf); err != nil {
		return err
	}
	s.data = buf.Bytes()
	return nil
}
