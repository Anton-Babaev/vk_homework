// Command inventory — CLI для управления инвентарём.
//
// main только собирает зависимости и завершает процесс; вся логика
// живёт в пакетах internal/*.
package main

import (
	"os"

	"github.com/Anton-Babaev/vk_homework/lesson_6/task1-inventory-cli/internal/cli"
	"github.com/Anton-Babaev/vk_homework/lesson_6/task1-inventory-cli/internal/storage"
)

func main() {
	app := cli.New(newStorage, os.Stdout, os.Stderr)
	os.Exit(app.Run(os.Args[1:]))
}

// newStorage — фабрика хранилища: по флагам команды выбирает файловое
// или in-memory хранилище. Это вся «сборка зависимостей», которая
// остаётся в main; команды работают только с интерфейсом cli.Storage.
func newStorage(path string, memory bool) cli.Storage {
	if memory {
		return storage.NewMemoryStorage()
	}
	return storage.NewFileStorage(path)
}
