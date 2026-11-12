package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// ============================================================
// 1. Ошибки предметов
// ============================================================

var (
	ErrBroken    = errors.New("item is broken")
	ErrExhausted = errors.New("item is exhausted")
)

// ============================================================
// 2. Типизированная ошибка парсинга
// ============================================================

type ParseItemError struct {
	Line  int
	Value string
	Err   error
}

func (e *ParseItemError) Error() string {
	return fmt.Sprintf("строка %d: не удалось разобрать %q: %v", e.Line, e.Value, e.Err)
}

func (e *ParseItemError) Unwrap() error { return e.Err }

// ============================================================
// 3. Интерфейс Item
// ============================================================

type Item interface {
	Name() string
	Weight() float64
	ID() int
	Use() (string, error)
	SetID(int)
}

var (
	_ Item = (*Weapon)(nil)
	_ Item = (*Armor)(nil)
	_ Item = (*Potion)(nil)
)

// ============================================================
// 4. Доменные типы
// ============================================================

type Weapon struct {
	id         int
	name       string
	weight     float64
	Damage     int
	Durability int
}

func (w *Weapon) Name() string    { return w.name }
func (w *Weapon) Weight() float64 { return w.weight }
func (w *Weapon) ID() int         { return w.id }
func (w *Weapon) SetID(id int)    { w.id = id }

func (w *Weapon) Use() (string, error) {
	if w.Durability <= 0 {
		return "", fmt.Errorf("оружие %q сломано: %w", w.name, ErrBroken)
	}
	w.Durability--
	return fmt.Sprintf("Взял %s (урон %d, прочность %d)", w.name, w.Damage, w.Durability), nil
}

type Armor struct {
	id      int
	name    string
	weight  float64
	Defense int
}

func (a *Armor) Name() string    { return a.name }
func (a *Armor) Weight() float64 { return a.weight }
func (a *Armor) ID() int         { return a.id }
func (a *Armor) SetID(id int)    { a.id = id }

func (a *Armor) Use() (string, error) {
	return fmt.Sprintf("Надел %s (защита %d)", a.name, a.Defense), nil
}

type Potion struct {
	id      int
	name    string
	weight  float64
	Effect  string
	Charges int
}

func (p *Potion) Name() string    { return p.name }
func (p *Potion) Weight() float64 { return p.weight }
func (p *Potion) ID() int         { return p.id }
func (p *Potion) SetID(id int)    { p.id = id }

func (p *Potion) Use() (string, error) {
	if p.Charges <= 0 {
		return "", fmt.Errorf("зелье %q исчерпано: %w", p.name, ErrExhausted)
	}
	p.Charges--
	return fmt.Sprintf("Выпил %s (эффект: %s, осталось зарядов: %d)", p.name, p.Effect, p.Charges), nil
}

// ============================================================
// 5. DescribeItem
// ============================================================

func DescribeItem(i Item) string {
	if i == nil {
		return "Предмет отсутствует"
	}
	switch v := i.(type) {
	case *Weapon:
		if v == nil {
			return "Предмет отсутствует"
		}
	case *Armor:
		if v == nil {
			return "Предмет отсутствует"
		}
	case *Potion:
		if v == nil {
			return "Предмет отсутствует"
		}
	}
	return fmt.Sprintf("%s (вес: %.2f)", i.Name(), i.Weight())
}

// ============================================================
// 6. Обобщённые функции
// ============================================================

func Filter[T any](items []T, predicate func(T) bool) []T {
	result := make([]T, 0, len(items))
	for _, item := range items {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

func Map[T any, R any](items []T, transform func(T) R) []R {
	result := make([]R, 0, len(items))
	for _, item := range items {
		result = append(result, transform(item))
	}
	return result
}

func Find[T any](items []T, condition func(T) bool) (T, bool) {
	for _, item := range items {
		if condition(item) {
			return item, true
		}
	}
	var zero T
	return zero, false
}

// ============================================================
// 7. Инвентарь
// ============================================================

type Inventory struct {
	Items  []Item
	nextID int
}

func (inv *Inventory) AddItem(item Item) int {
	inv.nextID++
	item.SetID(inv.nextID)
	inv.Items = append(inv.Items, item)
	return inv.nextID
}

func (inv *Inventory) FindByID(id int) (Item, bool) {
	for _, item := range inv.Items {
		if item.ID() == id {
			return item, true
		}
	}
	return nil, false
}

func (inv *Inventory) DeleteByID(id int) bool {
	for i, item := range inv.Items {
		if item.ID() == id {
			inv.Items = append(inv.Items[:i], inv.Items[i+1:]...)
			return true
		}
	}
	return false
}

func (inv *Inventory) Weapons() []*Weapon {
	filtered := Filter(inv.Items, func(i Item) bool {
		_, ok := i.(*Weapon)
		return ok
	})
	return Map(filtered, func(i Item) *Weapon { return i.(*Weapon) })
}

func (inv *Inventory) BrokenItems() []Item {
	return Filter(inv.Items, func(i Item) bool {
		switch v := i.(type) {
		case *Weapon:
			return v.Durability <= 0
		case *Potion:
			return v.Charges <= 0
		}
		return false
	})
}

func (inv *Inventory) ItemNames() []string {
	return Map(inv.Items, func(i Item) string { return i.Name() })
}

func (inv *Inventory) FindItemByName(name string) (Item, bool) {
	return Find(inv.Items, func(i Item) bool { return i.Name() == name })
}

// ============================================================
// 8. Сериализация
// ============================================================

func serialize(w io.Writer, item Item) error {
	switch v := item.(type) {
	case *Weapon:
		_, err := fmt.Fprintf(w, "Weapon|%d|%s|%d|%d|%.2f\n",
			v.id, v.name, v.Damage, v.Durability, v.weight)
		return err
	case *Armor:
		_, err := fmt.Fprintf(w, "Armor|%d|%s|%d|%.2f\n",
			v.id, v.name, v.Defense, v.weight)
		return err
	case *Potion:
		_, err := fmt.Fprintf(w, "Potion|%d|%s|%s|%d|%.2f\n",
			v.id, v.name, v.Effect, v.Charges, v.weight)
		return err
	default:
		return fmt.Errorf("неизвестный тип: %T", item)
	}
}

func parseItemLine(line string) (Item, error) {
	parts := strings.Split(line, "|")
	if len(parts) < 2 {
		return nil, fmt.Errorf("ожидалось минимум 2 поля, получено %d", len(parts))
	}
	switch parts[0] {
	case "Weapon":
		if len(parts) != 6 {
			return nil, fmt.Errorf("Weapon: ожидалось 6 полей, получено %d", len(parts))
		}
		id, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("Weapon id: %w", err)
		}
		damage, err := strconv.Atoi(parts[3])
		if err != nil {
			return nil, fmt.Errorf("Weapon Damage: %w", err)
		}
		dur, err := strconv.Atoi(parts[4])
		if err != nil {
			return nil, fmt.Errorf("Weapon Durability: %w", err)
		}
		weight, err := strconv.ParseFloat(parts[5], 64)
		if err != nil {
			return nil, fmt.Errorf("Weapon weight: %w", err)
		}
		return &Weapon{id: id, name: parts[2], Damage: damage, Durability: dur, weight: weight}, nil
	case "Armor":
		if len(parts) != 5 {
			return nil, fmt.Errorf("Armor: ожидалось 5 полей, получено %d", len(parts))
		}
		id, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("Armor id: %w", err)
		}
		defense, err := strconv.Atoi(parts[3])
		if err != nil {
			return nil, fmt.Errorf("Armor Defense: %w", err)
		}
		weight, err := strconv.ParseFloat(parts[4], 64)
		if err != nil {
			return nil, fmt.Errorf("Armor weight: %w", err)
		}
		return &Armor{id: id, name: parts[2], Defense: defense, weight: weight}, nil
	case "Potion":
		if len(parts) != 6 {
			return nil, fmt.Errorf("Potion: ожидалось 6 полей, получено %d", len(parts))
		}
		id, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("Potion id: %w", err)
		}
		charges, err := strconv.Atoi(parts[4])
		if err != nil {
			return nil, fmt.Errorf("Potion Charges: %w", err)
		}
		weight, err := strconv.ParseFloat(parts[5], 64)
		if err != nil {
			return nil, fmt.Errorf("Potion weight: %w", err)
		}
		return &Potion{id: id, name: parts[2], Effect: parts[3], Charges: charges, weight: weight}, nil
	default:
		return nil, fmt.Errorf("неизвестный тип: %q", parts[0])
	}
}

// ============================================================
// 9. Сохранение / загрузка
// ============================================================

func (inv *Inventory) Load(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	var errs []error
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		item, err := parseItemLine(line)
		if err != nil {
			errs = append(errs, &ParseItemError{Line: lineNo, Value: line, Err: err})
			continue
		}
		inv.Items = append(inv.Items, item)
		if item.ID() > inv.nextID {
			inv.nextID = item.ID()
		}
	}
	if err := scanner.Err(); err != nil {
		errs = append(errs, fmt.Errorf("чтение файла: %w", err))
	}
	return errors.Join(errs...)
}

func (inv *Inventory) Save(w io.Writer) error {
	bw := bufio.NewWriter(w)
	for _, item := range inv.Items {
		if err := serialize(bw, item); err != nil {
			return fmt.Errorf("сериализация %q: %w", item.Name(), err)
		}
	}
	if err := bw.Flush(); err != nil {
		return fmt.Errorf("flush: %w", err)
	}
	return nil
}

func loadInventoryFromFile(path string) (*Inventory, error) {
	inv := &Inventory{}
	f, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return inv, nil // файла нет — пустой инвентарь, не ошибка
		}
		return nil, fmt.Errorf("открытие %q: %w", path, err)
	}
	defer f.Close()
	if err := inv.Load(f); err != nil {
		// Явно извлекаем ParseItemError через errors.As
		var pie *ParseItemError
		if errors.As(err, &pie) {
			return inv, fmt.Errorf("битая строка %d: %q (причина: %v): %w",
				pie.Line, pie.Value, pie.Err, err)
		}
		return inv, fmt.Errorf("загрузка %q: %w", path, err)
	}
	return inv, nil
}

func saveInventoryToFile(path string, inv *Inventory) (err error) {
	dir := filepath.Dir(path)
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
	if err := os.Rename(tmpName, path); err != nil {
		return fmt.Errorf("rename: %w", err)
	}
	return nil
}

// ============================================================
// 10. panic / recover
// ============================================================

// UseDangerousItem демонстрирует преобразование паники в ошибку.
//
// Ответы на вопросы задания:
//
//	1.Ошибка сломанного предмета — не паника: это штатное доменное
//	  состояние, которое вызывающий может обработать. Panic
//	  в Go — для программных ошибок, и он обрушил бы весь CLI.
//	2.После panic(...) не выполнится return item.Use() — управление
//	  сразу уйдёт в отложенные вызовы (defer), где сработает recover.
//	3.Recover работает только во время раскрутки стека после panic,
//	  поэтому он имеет смысл только внутри defer-функции.
func UseDangerousItem(item Item) (message string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("паника при использовании предмета: %v", r)
		}
	}()
	if w, ok := item.(*Weapon); ok && w.name == "Проклятый клинок" {
		panic("проклятый клинок вырвался из рук")
	}
	return item.Use()
}

// ============================================================
// 11. CLI
// ============================================================

const defaultInventoryPath = "inventory.txt"

func main() {
	code := run(os.Args[1:], os.Stdout, os.Stderr)
	os.Exit(code)
}

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) < 1 {
		fmt.Fprintln(stderr, "usage: inventory <list|add|use|delete> [flags]")
		return 2
	}
	switch args[0] {
	case "list":
		return cmdList(args[1:], stdout, stderr)
	case "add":
		return cmdAdd(args[1:], stdout, stderr)
	case "use":
		return cmdUse(args[1:], stdout, stderr)
	case "delete":
		return cmdDelete(args[1:], stdout, stderr)
	default:
		fmt.Fprintf(stderr, "неизвестная команда: %q\n", args[0])
		return 2
	}
}

func cmdList(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("list", flag.ContinueOnError)
	fs.SetOutput(stderr)
	path := fs.String("file", defaultInventoryPath, "путь к файлу инвентаря")
	if err := fs.Parse(args); err != nil {
		return 2
	}

	inv, err := loadInventoryFromFile(*path)
	if err != nil {
		fmt.Fprintf(stderr, "загрузка: %v\n", err)
		if inv == nil {
			return 1
		}
	}
	if len(inv.Items) == 0 {
		fmt.Fprintln(stdout, "инвентарь пуст")
		return 0
	}
	for _, item := range inv.Items {
		fmt.Fprintf(stdout, "ID=%d | %T | %s (вес: %.2f)\n",
			item.ID(), item, item.Name(), item.Weight())
	}
	return 0
}

func cmdAdd(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("add", flag.ContinueOnError)
	fs.SetOutput(stderr)
	path := fs.String("file", defaultInventoryPath, "путь к файлу инвентаря")
	kind := fs.String("type", "weapon", "тип: weapon|armor|potion")
	name := fs.String("name", "", "название предмета")
	weight := fs.Float64("weight", 0, "вес предмета")
	damage := fs.Int("damage", 0, "урон (weapon)")
	durability := fs.Int("durability", 1, "прочность (weapon)")
	defense := fs.Int("defense", 0, "защита (armor)")
	effect := fs.String("effect", "", "эффект (potion)")
	charges := fs.Int("charges", 1, "заряды (potion)")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *name == "" {
		fmt.Fprintln(stderr, "флаг -name обязателен")
		return 2
	}
	if *weight <= 0 {
		fmt.Fprintln(stderr, "флаг -weight должен быть > 0")
		return 2
	}

	inv, loadErr := loadInventoryFromFile(*path)
	if inv == nil {
		fmt.Fprintf(stderr, "загрузка: %v\n", loadErr)
		return 1
	}
	if loadErr != nil {
		fmt.Fprintf(stderr, "предупреждение: %v\n", loadErr)
	}

	var item Item
	switch *kind {
	case "weapon":
		item = &Weapon{name: *name, weight: *weight, Damage: *damage, Durability: *durability}
	case "armor":
		item = &Armor{name: *name, weight: *weight, Defense: *defense}
	case "potion":
		item = &Potion{name: *name, weight: *weight, Effect: *effect, Charges: *charges}
	default:
		fmt.Fprintf(stderr, "неизвестный тип: %q\n", *kind)
		return 2
	}

	id := inv.AddItem(item)
	if err := saveInventoryToFile(*path, inv); err != nil {
		fmt.Fprintf(stderr, "сохранение: %v\n", err)
		return 1
	}
	fmt.Fprintf(stdout, "добавлено: ID=%d %s\n", id, item.Name())
	return 0
}

func cmdUse(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("use", flag.ContinueOnError)
	fs.SetOutput(stderr)
	path := fs.String("file", defaultInventoryPath, "путь к файлу инвентаря")
	id := fs.Int("id", 0, "ID предмета")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *id <= 0 {
		fmt.Fprintln(stderr, "флаг -id обязателен")
		return 2
	}

	inv, loadErr := loadInventoryFromFile(*path)
	if inv == nil {
		fmt.Fprintf(stderr, "загрузка: %v\n", loadErr)
		return 1
	}
	if loadErr != nil {
		fmt.Fprintf(stderr, "предупреждение: %v\n", loadErr)
	}

	item, ok := inv.FindByID(*id)
	if !ok {
		fmt.Fprintf(stderr, "предмет с id=%d не найден\n", *id)
		return 1
	}

	msg, err := UseDangerousItem(item)
	if err != nil {
		switch {
		case errors.Is(err, ErrBroken):
			fmt.Fprintf(stderr, "нельзя использовать, предмет сломан: %v\n", err)
		case errors.Is(err, ErrExhausted):
			fmt.Fprintf(stderr, "нельзя использовать, заряды кончились: %v\n", err)
		default:
			fmt.Fprintf(stderr, "ошибка использования: %v\n", err)
		}
		return 1
	}
	fmt.Fprintln(stdout, msg)

	if err := saveInventoryToFile(*path, inv); err != nil {
		fmt.Fprintf(stderr, "сохранение: %v\n", err)
		return 1
	}
	return 0
}

func cmdDelete(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("delete", flag.ContinueOnError)
	fs.SetOutput(stderr)
	path := fs.String("file", defaultInventoryPath, "путь к файлу инвентаря")
	id := fs.Int("id", 0, "ID предмета")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *id <= 0 {
		fmt.Fprintln(stderr, "флаг -id обязателен")
		return 2
	}

	inv, loadErr := loadInventoryFromFile(*path)
	if inv == nil {
		fmt.Fprintf(stderr, "загрузка: %v\n", loadErr)
		return 1
	}
	if loadErr != nil {
		fmt.Fprintf(stderr, "предупреждение: %v\n", loadErr)
	}

	if !inv.DeleteByID(*id) {
		fmt.Fprintf(stderr, "предмет с id=%d не найден\n", *id)
		return 1
	}
	if err := saveInventoryToFile(*path, inv); err != nil {
		fmt.Fprintf(stderr, "сохранение: %v\n", err)
		return 1
	}
	fmt.Fprintf(stdout, "удалено: ID=%d\n", *id)
	return 0
}
