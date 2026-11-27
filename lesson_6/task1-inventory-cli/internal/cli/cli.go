// Package cli реализует разбор команд и их выполнение. Пакет зависит
// от inventory и от абстрактного хранилища (интерфейс Storage объявлен
// здесь, со стороны потребителя).
package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/Anton-Babaev/vk_homework/lesson_6/task1-inventory-cli/internal/inventory"
	"github.com/Anton-Babaev/vk_homework/lesson_6/task1-inventory-cli/internal/item"
)

// Storage — то, что нужно CLI от хранилища. Объявлен со стороны
// потребителя: любая реализация (файл, память) подходит без изменений
// вызывающего кода.
type Storage interface {
	Load() (*inventory.Inventory, error)
	Save(*inventory.Inventory) error
}

// StorageFactory создаёт хранилище по флагам команды. main передаёт
// фабрику, поэтому выбор «файл или память» и путь остаются на стороне
// сборки зависимостей, а команды не зависят от конкретной реализации.
type StorageFactory func(path string, memory bool) Storage

// App — приложение CLI с внедрёнными зависимостями.
type App struct {
	newStorage StorageFactory
	stdout     io.Writer
	stderr     io.Writer
}

// New создаёт приложение с фабрикой хранилища и потоками вывода.
func New(newStorage StorageFactory, stdout, stderr io.Writer) *App {
	return &App{newStorage: newStorage, stdout: stdout, stderr: stderr}
}

// defaultInventoryPath — файл инвентаря по умолчанию.
const defaultInventoryPath = "inventory.txt"

// Run разбирает имя команды и выполняет её. Возвращает код выхода.
func (a *App) Run(args []string) int {
	if len(args) < 1 {
		fmt.Fprintln(a.stderr, "usage: inventory <list|add|use|delete> [flags]")
		return 2
	}
	switch args[0] {
	case "list":
		return a.cmdList(args[1:])
	case "add":
		return a.cmdAdd(args[1:])
	case "use":
		return a.cmdUse(args[1:])
	case "delete":
		return a.cmdDelete(args[1:])
	default:
		fmt.Fprintf(a.stderr, "неизвестная команда: %q\n", args[0])
		return 2
	}
}

// newFlagSet создаёт FlagSet с общим для всех команд поведением.
func (a *App) newFlagSet(name string) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(a.stderr)
	return fs
}

// parseFlags разбирает флаги и проверяет отсутствие лишних аргументов.
// Возвращает код выхода и признак успеха.
func (a *App) parseFlags(fs *flag.FlagSet, args []string) (int, bool) {
	if err := fs.Parse(args); err != nil {
		return 2, false
	}
	if fs.NArg() > 0 {
		fmt.Fprintf(a.stderr, "лишние аргументы: %v\n", fs.Args())
		return 2, false
	}
	return 0, true
}

// storageFlags регистрирует общие флаги хранилища и возвращает указатели
// на их значения. -file и -memory — обычные флаги подкоманды, поэтому
// работают там же, где и в ДЗ5: после имени команды. Значения заполняются
// при разборе аргументов, поэтому читать их нужно после parseFlags.
func (a *App) storageFlags(fs *flag.FlagSet) (*string, *bool) {
	path := fs.String("file", defaultInventoryPath, "путь к файлу инвентаря")
	memory := fs.Bool("memory", false, "хранить в памяти (на время одного запуска)")
	return path, memory
}

func (a *App) cmdList(args []string) int {
	fs := a.newFlagSet("list")
	path, memory := a.storageFlags(fs)
	if code, ok := a.parseFlags(fs, args); !ok {
		return code
	}

	store := a.newStorage(*path, *memory)
	inv, err := store.Load()
	if err != nil {
		fmt.Fprintf(a.stderr, "загрузка: %v\n", err)
		return 1
	}
	if inv.Len() == 0 {
		fmt.Fprintln(a.stdout, "инвентарь пуст")
		return 0
	}
	for _, it := range inv.Items() {
		fmt.Fprintf(a.stdout, "ID=%d | %s | %s (вес: %.2f)\n",
			it.ID(), item.KindName(it), it.Name(), it.Weight())
	}
	return 0
}

func (a *App) cmdAdd(args []string) int {
	fs := a.newFlagSet("add")
	path, memory := a.storageFlags(fs)
	kind := fs.String("type", "weapon", "тип: weapon|armor|potion")
	name := fs.String("name", "", "название предмета")
	weight := fs.Float64("weight", 0, "вес предмета")
	damage := fs.Int("damage", 0, "урон (weapon)")
	durability := fs.Int("durability", 1, "прочность (weapon/armor)")
	defense := fs.Int("defense", 0, "защита (armor)")
	effect := fs.String("effect", "", "эффект (potion)")
	charges := fs.Int("charges", 1, "заряды (potion)")
	if code, ok := a.parseFlags(fs, args); !ok {
		return code
	}

	if *name == "" {
		fmt.Fprintln(a.stderr, "флаг -name обязателен")
		return 2
	}
	if *weight <= 0 {
		fmt.Fprintln(a.stderr, "флаг -weight должен быть > 0")
		return 2
	}
	// Запрещаем разделители и переводы строк: иначе запись станет
	// нечитаемой при следующей загрузке.
	if strings.ContainsAny(*name, "|\n\r") {
		fmt.Fprintln(a.stderr, "имя не должно содержать | или перевод строки")
		return 2
	}
	if *kind == "potion" && strings.ContainsAny(*effect, "|\n\r") {
		fmt.Fprintln(a.stderr, "эффект не должен содержать | или перевод строки")
		return 2
	}

	store := a.newStorage(*path, *memory)
	inv, loadErr := store.Load()
	if loadErr != nil {
		// Изменяющая команда не должна перезаписывать данные, если в них
		// есть повреждённые записи — иначе они будут потеряны.
		fmt.Fprintf(a.stderr, "загрузка: %v\n", loadErr)
		return 1
	}

	var it item.Item
	switch *kind {
	case "weapon":
		it = item.NewWeapon(*name, *weight, *damage, *durability)
	case "armor":
		it = item.NewArmor(*name, *weight, *defense, *durability)
	case "potion":
		it = item.NewPotion(*name, *weight, *effect, *charges)
	default:
		fmt.Fprintf(a.stderr, "неизвестный тип: %q\n", *kind)
		return 2
	}

	id := inv.AddItem(it)
	if err := store.Save(inv); err != nil {
		fmt.Fprintf(a.stderr, "сохранение: %v\n", err)
		return 1
	}
	fmt.Fprintf(a.stdout, "добавлено: ID=%d %s\n", id, it.Name())
	return 0
}

func (a *App) cmdUse(args []string) int {
	fs := a.newFlagSet("use")
	path, memory := a.storageFlags(fs)
	id := fs.Int("id", 0, "ID предмета")
	if code, ok := a.parseFlags(fs, args); !ok {
		return code
	}
	if *id <= 0 {
		fmt.Fprintln(a.stderr, "флаг -id обязателен")
		return 2
	}

	store := a.newStorage(*path, *memory)
	inv, loadErr := store.Load()
	if loadErr != nil {
		fmt.Fprintf(a.stderr, "загрузка: %v\n", loadErr)
		return 1
	}

	it, ok := inv.FindByID(*id)
	if !ok {
		fmt.Fprintf(a.stderr, "предмет с id=%d не найден\n", *id)
		return 1
	}

	msg, err := useDangerousItem(it)
	if err != nil {
		switch {
		case errors.Is(err, item.ErrBroken):
			fmt.Fprintf(a.stderr, "нельзя использовать, предмет сломан: %v\n", err)
		case errors.Is(err, item.ErrExhausted):
			fmt.Fprintf(a.stderr, "нельзя использовать, заряды кончились: %v\n", err)
		default:
			fmt.Fprintf(a.stderr, "ошибка использования: %v\n", err)
		}
		return 1
	}

	// Сначала сохраняем, и только при успехе сообщаем об успехе.
	if err := store.Save(inv); err != nil {
		fmt.Fprintf(a.stderr, "сохранение: %v\n", err)
		return 1
	}
	fmt.Fprintln(a.stdout, msg)
	return 0
}

func (a *App) cmdDelete(args []string) int {
	fs := a.newFlagSet("delete")
	path, memory := a.storageFlags(fs)
	id := fs.Int("id", 0, "ID предмета")
	if code, ok := a.parseFlags(fs, args); !ok {
		return code
	}
	if *id <= 0 {
		fmt.Fprintln(a.stderr, "флаг -id обязателен")
		return 2
	}

	store := a.newStorage(*path, *memory)
	inv, loadErr := store.Load()
	if loadErr != nil {
		fmt.Fprintf(a.stderr, "загрузка: %v\n", loadErr)
		return 1
	}

	if !inv.DeleteByID(*id) {
		fmt.Fprintf(a.stderr, "предмет с id=%d не найден\n", *id)
		return 1
	}
	if err := store.Save(inv); err != nil {
		fmt.Fprintf(a.stderr, "сохранение: %v\n", err)
		return 1
	}
	fmt.Fprintf(a.stdout, "удалено: ID=%d\n", *id)
	return 0
}

// useDangerousItem демонстрирует преобразование паники в ошибку.
//
// Ответы на вопросы задания:
//
//  1. Ошибка сломанного предмета — не паника: это штатное доменное
//     состояние, которое вызывающий может обработать. Panic
//     в Go — для программных ошибок, и он обрушил бы весь CLI.
//  2. После panic(...) не выполнится return item.Use() — управление
//     сразу уйдёт в отложенные вызовы (defer), где сработает recover.
//  3. Recover работает только во время раскрутки стека после panic,
//     поэтому он имеет смысл только внутри defer-функции.
func useDangerousItem(it item.Item) (message string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("паника при использовании предмета: %v", r)
		}
	}()
	if w, ok := it.(*item.Weapon); ok && w.Name() == "Проклятый клинок" {
		panic("проклятый клинок вырвался из рук")
	}
	return it.Use()
}
