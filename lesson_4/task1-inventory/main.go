package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// TODO: Реализуйте структуры:
// TODO: - Weapon: name (string), weight (float64), Damage (int), Durability (int)
// TODO: - Armor: name (string), weight (float64), Defense (int)
// TODO: - Potion: name (string), weight (float64), Effect (string), Charges (int)

type Weapon struct {
	name       string
	weight     float64
	Damage     int
	Durability int
}

type Armor struct {
	name    string
	weight  float64
	Defense int
}

type Potion struct {
	name    string
	weight  float64
	Effect  string
	Charges int
}

// TODO: Можете добавить свои структуры и/или немного модифицировать эти :)
// TODO:
// TODO: Для каждой структуры реализуйте:
// TODO: - Метод Use() string (описание использования, например "Используется <имя>" и т.д.)

func (w *Weapon) Name() string    { return w.name }
func (w *Weapon) Weight() float64 { return w.weight }
func (w *Weapon) Use() string {
	return fmt.Sprintf("В левую руку взял %s (урон %d)", w.name, w.Damage)
}

func (a *Armor) Name() string    { return a.name }
func (a *Armor) Weight() float64 { return a.weight }
func (a *Armor) Use() string {
	return fmt.Sprintf("В правой руке сжал %s (защита %d)", a.name, a.Defense)
}

func (p *Potion) Name() string    { return p.name }
func (p *Potion) Weight() float64 { return p.weight }
func (p *Potion) Use() string {
	return fmt.Sprintf("С размаху влупил %s (%s, до похмелья %d часа)", p.name, p.Effect, p.Charges)
}

// TODO: - Методы интерфейса Item

type Item interface {
	Name() string
	Weight() float64
}

// TODO: Добавьте compile-time проверки, что *Weapon, *Armor и *Potion реализуют Item:
// TODO: var _ Item = (*Weapon)(nil)
// не работает без:
// func (w *Weapon) Name() string    { return w.name }
// func (w *Weapon) Weight() float64 { return w.weight }
var (
	_ Item = (*Weapon)(nil)
	_ Item = (*Potion)(nil)
	_ Item = (*Armor)(nil)
)

// TODO: Реализуйте функцию

func DescribeItem(i Item) string {
	if i == nil {
		return "Предмет отсутствует"
	}
	// Функция должна возвращать:
	// - "Предмет отсутствует" если i == nil или внутри nil-указатель (var weapon *Weapon)
	// - "<название> (вес: <вес>)" в остальных случаях
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

func Filter[T any](items []T, predicate func(T) bool) []T {
	// TODO: Верните новый слайс только с элементами, для которых predicate вернул true
	result := make([]T, 0, len(items))
	for _, item := range items {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

func Map[T any, R any](items []T, transform func(T) R) []R {
	// TODO: Примените transform к каждому элементу и верните слайс с результатами
	result := make([]R, 0, len(items))
	for _, item := range items {
		result = append(result, transform(item))
	}
	return result
}

func Find[T any](items []T, condition func(T) bool) (T, bool) {
	// TODO: Найдите первый элемент, удовлетворяющий condition и верните элемент и true или zero value и false
	for _, item := range items {
		if condition(item) {
			return item, true
		}
	}
	var zero T
	return zero, false
}

type Inventory struct {
	Items []Item
}

func (inv *Inventory) AddItem(item Item) {
	inv.Items = append(inv.Items, item)
}

func (inv *Inventory) Weapons() []*Weapon {
	// TODO: Используйте Filter для отбора Weapon, затем Map для преобразования Item -> *Weapon
	filtered := Filter(inv.Items, func(i Item) bool {
		_, ok := i.(*Weapon)
		return ok
	})
	return Map(filtered, func(i Item) *Weapon {
		return i.(*Weapon)
	})
}

func (inv *Inventory) BrokenItems() []Item {
	// TODO: Используйте Filter для отбора:
	// TODO: - Weapon: Durability <= 0
	// TODO: - Potion: Charges <= 0
	// TODO:
	// TODO: Подсказка: поможет приведение типов - item.(type)
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
	// TODO: Используйте Map для преобразования []Item -> []string
	return Map(inv.Items, func(i Item) string {
		return i.Name()
	})
}

func (inv *Inventory) FindItemByName(name string) (Item, bool) {
	// TODO: Используйте Find для поиска по имени
	return Find(inv.Items, func(i Item) bool {
		return i.Name() == name
	})
}

// TODO: Бонус: реализуйте интефейс Storable для Weapon и Armor:
// TODO: - Weapon: формат "Weapon|Name|Damage|Durability|Weight"
// TODO: - Armor: формат "Armor|Name|Defense|Weight"

type Storable interface {
	Serialize(w io.Writer)
	Deserialize(r io.Reader)
}

func (w *Weapon) Serialize(out io.Writer) {
	fmt.Fprintf(out, "Оружие|%s|%d|%d|%.2f\n", w.name, w.Damage, w.Durability, w.weight)
}

func (w *Weapon) Deserialize(in io.Reader) {
	line, _ := bufio.NewReader(in).ReadString('\n')
	parts := strings.Split(strings.TrimSpace(line), "|")
	if len(parts) == 5 {
		w.name = parts[1]
		fmt.Sscanf(parts[2], "%d", &w.Damage)
		fmt.Sscanf(parts[3], "%d", &w.Durability)
		fmt.Sscanf(parts[4], "%f", &w.weight)
	}
}

func (a *Armor) Serialize(out io.Writer) {
	fmt.Fprintf(out, "Броня|%s|%d|%.2f\n", a.name, a.Defense, a.weight)
}

func (a *Armor) Deserialize(in io.Reader) {
	line, _ := bufio.NewReader(in).ReadString('\n')
	parts := strings.Split(strings.TrimSpace(line), "|")
	if len(parts) == 4 {
		a.name = parts[1]
		fmt.Sscanf(parts[2], "%d", &a.Defense)
		fmt.Sscanf(parts[3], "%f", &a.weight)
	}
}

func (inv *Inventory) Save(w io.Writer) {
	// TODO: Бонус: сделайте сохранение/загрузку инвентаря в/из файла
	for _, item := range inv.Items {
		if s, ok := item.(Storable); ok {
			s.Serialize(w)
		}
	}
}

func (inv *Inventory) Load(r io.Reader) {
	// TODO: Бонус: сделайте сохранение/загрузку инвентаря в/из файла
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "|")
		switch parts[0] {
		case "Оружие":
			w := &Weapon{}
			w.Deserialize(strings.NewReader(line + "\n"))
			inv.AddItem(w)
		case "Броня":
			a := &Armor{}
			a.Deserialize(strings.NewReader(line + "\n"))
			inv.AddItem(a)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Ошибка при чтении инвентаря: %v\n", err)
	}
}

func main() {
	// TODO: Создайте инвентарь и добавьте:
	inv := &Inventory{}

	// TODO: - Оружие: "Меч" (урон 10, прочность 5)
	inv.AddItem(&Weapon{name: "Меч", weight: 3.5, Damage: 10, Durability: 5})

	// TODO: - Броню: "Щит" (защита 5, вес 4.5)
	inv.AddItem(&Armor{name: "Щит", weight: 4.5, Defense: 5})

	// TODO: - Зелье: "Лечебное" (+50 HP, 3 заряда)
	inv.AddItem(&Potion{name: "Эль ядрёный", weight: 0.5, Effect: "+50 HP", Charges: 3})

	// TODO: - Оружие: "Сломанный лук" (урон 5, прочность 0)
	inv.AddItem(&Weapon{name: "Сломанный лук", weight: 2.0, Damage: 5, Durability: 0})

	// TODO: Реализуйте логику/вызовы:
	// TODO: 1. Use предмета с выводом в консоль
	fmt.Println("=== Действие ===")
	for _, item := range inv.Items {
		if u, ok := item.(interface{ Use() string }); ok {
			fmt.Println(u.Use())
		}
	}

	// TODO: 2. DescribeItem с предметом, с nil и с nil-указателем (var weapon *Weapon), так же с выводом в консоль
	fmt.Println("\n=== Описание ===")
	fmt.Println(DescribeItem(inv.Items[0]))
	fmt.Println(DescribeItem(nil))
	var weapon *Weapon
	fmt.Println(DescribeItem(weapon))

	// TODO: 3. Вывести в консоль результат вызова Weapons (должны вернуться только меч и лук)
	fmt.Println("\n=== Оружие ===")
	for _, w := range inv.Weapons() {
		fmt.Printf("%s (урон %d, прочность %d)\n", w.name, w.Damage, w.Durability)
	}

	// TODO: 4. Вывести в консоль результат вызова BrokenItems (должен вернуть сломанный лук)
	fmt.Println("\n=== Сломанные предметы ===")
	for _, b := range inv.BrokenItems() {
		fmt.Println(b.Name())
	}

	// TODO: 5. Вывести в консоль результат вызова ItemNames (все названия)
	fmt.Println("\n=== Все предметы ===")
	fmt.Println(strings.Join(inv.ItemNames(), ", "))

	// TODO: 6. Вывести в консоль результат вызова FindItemByName (поиск "Щит")
	fmt.Println("\n=== Поиск по названию ===")
	if item, ok := inv.FindItemByName("Щит"); ok {
		fmt.Println(DescribeItem(item))
	} else {
		fmt.Println("не найдено")
	}

	// TODO: Бонус: сделайте сохранение инвентаря в файл и загрузку инвентаря из файла
	fmt.Println("\n=== Сохранить/Загрузить ===")

	// --- Save в файл ---
	file, err := os.Create("inventory.txt")
	if err != nil {
		fmt.Println("Ошибка создания файла:", err)
		return
	}
	defer file.Close()

	inv.Save(file)
	fmt.Println("Сохранено в inventory.txt")

	// --- Load из файла ---
	file, err = os.Open("inventory.txt")
	if err != nil {
		fmt.Println("Ошибка открытия файла:", err)
		return
	}
	defer file.Close()

	newInv := &Inventory{}
	newInv.Load(file)
	fmt.Println("Загружено из inventory.txt:")
	fmt.Println(newInv.ItemNames())
}
