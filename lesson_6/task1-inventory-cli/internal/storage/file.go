// Package storage реализует хранилища инвентаря: файловое и in-memory.
// Обе реализации соблюдают одинаковые правила загрузки и сохранения,
// поэтому взаимозаменяемы для вызывающего кода.
package storage

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Anton-Babaev/vk_homework/lesson_6/task1-inventory-cli/internal/inventory"
)

// FileStorage хранит инвентарь в текстовом файле.
type FileStorage struct {
	path string
}

// NewFileStorage создаёт файловое хранилище по пути.
func NewFileStorage(path string) *FileStorage {
	return &FileStorage{path: path}
}

// Load читает инвентарь из файла. Отсутствие файла — не ошибка:
// возвращается пустой инвентарь.
func (s *FileStorage) Load() (*inventory.Inventory, error) {
	inv := inventory.New()
	f, err := os.Open(s.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return inv, nil
		}
		return nil, fmt.Errorf("открытие %q: %w", s.path, err)
	}
	defer f.Close()

	if err := inv.Load(f); err != nil {
		return inv, wrapLoadError(s.path, err)
	}
	return inv, nil
}

// Save атомарно записывает инвентарь в файл через временный файл и rename.
func (s *FileStorage) Save(inv *inventory.Inventory) (err error) {
	dir := filepath.Dir(s.path)
	tmp, err := os.CreateTemp(dir, ".inventory-*")
	if err != nil {
		return fmt.Errorf("create temp: %w", err)
	}
	tmpName := tmp.Name()
	defer func() {
		if err != nil {
			tmp.Close()
			os.Remove(tmpName)
		}
	}()

	if err := inv.Save(tmp); err != nil {
		return fmt.Errorf("запись: %w", err)
	}
	if err := tmp.Sync(); err != nil {
		return fmt.Errorf("sync: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close: %w", err)
	}
	if err := os.Rename(tmpName, s.path); err != nil {
		return fmt.Errorf("rename: %w", err)
	}
	return nil
}
