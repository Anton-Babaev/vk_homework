package storage

import (
	"errors"
	"fmt"

	"github.com/Anton-Babaev/vk_homework/lesson_6/task1-inventory-cli/internal/inventory"
)

// wrapLoadError приводит ошибку загрузки к единому виду для обеих
// реализаций хранилища (требование бонуса ДЗ6: одинаковые правила
// обработки ошибок). Типизированная ParseItemError оборачивается с
// номером и содержимым битой строки; прочие ошибки — с путём файла.
func wrapLoadError(path string, err error) error {
	var pie *inventory.ParseItemError
	if errors.As(err, &pie) {
		return fmt.Errorf("битая строка %d: %q (причина: %v): %w",
			pie.Line, pie.Value, pie.Err, err)
	}
	if path == "" {
		return err
	}
	return fmt.Errorf("загрузка %q: %w", path, err)
}
