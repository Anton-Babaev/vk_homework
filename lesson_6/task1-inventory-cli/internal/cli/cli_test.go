package cli_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/Anton-Babaev/vk_homework/lesson_6/task1-inventory-cli/internal/cli"
	"github.com/Anton-Babaev/vk_homework/lesson_6/task1-inventory-cli/internal/inventory"
	"github.com/Anton-Babaev/vk_homework/lesson_6/task1-inventory-cli/internal/item"
)

// fakeStorage — минимальная реализация cli.Storage для тестов.
type fakeStorage struct {
	inv *inventory.Inventory
	err error
}

func (f *fakeStorage) Load() (*inventory.Inventory, error) {
	if f.err != nil {
		return nil, f.err
	}
	if f.inv == nil {
		f.inv = inventory.New()
	}
	return f.inv, nil
}

func (f *fakeStorage) Save(inv *inventory.Inventory) error {
	f.inv = inv
	return f.err
}

// run собирает App с фабрикой, которая запоминает флаги хранилища.
func run(t *testing.T, args ...string) (code int, stdout, stderr, path string, memory bool) {
	t.Helper()
	var out, errB bytes.Buffer
	app := cli.New(func(p string, m bool) cli.Storage {
		path, memory = p, m
		return &fakeStorage{}
	}, &out, &errB)
	code = app.Run(args)
	return code, out.String(), errB.String(), path, memory
}

// TestList_FileFlagAfterCommand — регрессия: -file должен работать
// ПОСЛЕ имени команды, как в ДЗ5.
func TestList_FileFlagAfterCommand(t *testing.T) {
	code, stdout, _, path, memory := run(t, "list", "-file", "data.txt")
	if code != 0 {
		t.Fatalf("code = %d, want 0", code)
	}
	if path != "data.txt" {
		t.Fatalf("path = %q, want %q", path, "data.txt")
	}
	if memory {
		t.Fatal("memory не должен быть установлен")
	}
	if !strings.Contains(stdout, "инвентарь пуст") {
		t.Fatalf("stdout = %q, ожидался «инвентарь пуст»", stdout)
	}
}

// TestList_MemoryFlag — -memory передаётся в фабрику.
func TestList_MemoryFlag(t *testing.T) {
	_, _, _, _, memory := run(t, "list", "-memory")
	if !memory {
		t.Fatal("-memory не распознан")
	}
}

// TestFileFlagBeforeCommandIsUnknown — как в ДЗ5: флаг до команды —
// неизвестная команда, код 2.
func TestFileFlagBeforeCommandIsUnknown(t *testing.T) {
	code, _, stderr, _, _ := run(t, "-file", "data.txt", "list")
	if code != 2 {
		t.Fatalf("code = %d, want 2", code)
	}
	if !strings.Contains(stderr, "неизвестная команда") {
		t.Fatalf("stderr = %q, ожидалось «неизвестная команда»", stderr)
	}
}

// TestExitCodes — коды 0/1/2 для разных ситуаций.
func TestExitCodes(t *testing.T) {
	cases := []struct {
		name string
		args []string
		want int
	}{
		{"нет аргументов", nil, 2},
		{"неизвестная команда", []string{"foo"}, 2},
		{"неверный флаг", []string{"list", "-unknown"}, 2},
		{"лишние аргументы", []string{"list", "extra"}, 2},
		{"пустой list", []string{"list"}, 0},
		{"use без id", []string{"use"}, 2},
		{"delete без id", []string{"delete"}, 2},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			code, _, _, _, _ := run(t, tc.args...)
			if code != tc.want {
				t.Fatalf("code = %d, want %d", code, tc.want)
			}
		})
	}
}

// TestUse_PanicRecovered — «Проклятый клинок» превращает панику в ошибку.
func TestUse_PanicRecovered(t *testing.T) {
	inv := inventory.New()
	inv.AddItem(item.NewWeapon("Проклятый клинок", 1, 99, 5)) // ID=1

	var out, errB bytes.Buffer
	app := cli.New(func(string, bool) cli.Storage {
		return &fakeStorage{inv: inv}
	}, &out, &errB)

	code := app.Run([]string{"use", "-id", "1"})
	if code != 1 {
		t.Fatalf("code = %d, want 1", code)
	}
	if !strings.Contains(errB.String(), "паника при использовании") {
		t.Fatalf("stderr = %q, ожидалось сообщение о панике", errB.String())
	}
}

// TestAdd_SuccessToStdoutAndSaved — add печатает успех в stdout и вызывает Save.
func TestAdd_SuccessToStdoutAndSaved(t *testing.T) {
	store := &fakeStorage{inv: inventory.New()}
	var out, errB bytes.Buffer
	app := cli.New(func(string, bool) cli.Storage { return store }, &out, &errB)

	code := app.Run([]string{"add", "-name", "Меч", "-weight", "2.5", "-damage", "10"})
	if code != 0 {
		t.Fatalf("code = %d, want 0", code)
	}
	if !strings.Contains(out.String(), "добавлено: ID=1 Меч") {
		t.Fatalf("stdout = %q", out.String())
	}
	if store.inv.Len() != 1 {
		t.Fatalf("сохранено предметов %d, want 1", store.inv.Len())
	}
}
