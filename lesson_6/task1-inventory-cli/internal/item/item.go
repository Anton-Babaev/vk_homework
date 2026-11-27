// Package item содержит доменные типы предметов инвентаря и их
// сериализацию в текстовый формат. Пакет не знает ни о хранилище,
// ни о CLI — это чистые данные и правила одного предмета.
package item

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Ошибки предметов. Это штатные доменные состояния, а не паника:
// вызывающий может их обработать через errors.Is.
var (
	ErrBroken    = errors.New("item is broken")
	ErrExhausted = errors.New("item is exhausted")
)

// Item — общий контракт любого предмета инвентаря.
type Item interface {
	Name() string
	Weight() float64
	ID() int
	Use() (string, error)
	SetID(int)
}

// Storable — предмет, который умеет сохраняться и загружаться.
// Объявлен здесь, рядом с реализациями; потребитель (хранилище)
// использует его через этот же интерфейс.
type Storable interface {
	Serialize(w io.Writer) error
	Deserialize(r io.Reader) error
}

// ============================================================
// Общий разбор строки формата "Type|id|name|field3|field4|weight"
// ============================================================

// fields — разобранные поля записи. Все три типа предметов имеют
// одинаковую структуру из 6 полей, поэтому разбор общий, а
// специфика (что лежит в field3/field4) — в конкретном типе.
type fields struct {
	id     int
	name   string
	field3 string
	field4 string
	weight float64
}

// parseFields разбирает одну строку формата и валидирует общую часть:
// количество полей, id, weight. Специфичные поля (field3/field4)
// остаются строками — их интерпретирует конкретный тип.
func parseFields(line, kind string) (fields, error) {
	parts := strings.Split(strings.TrimSpace(line), "|")
	if len(parts) != 6 {
		return fields{}, fmt.Errorf("%s: ожидалось 6 полей, получено %d", kind, len(parts))
	}
	id, err := strconv.Atoi(parts[1])
	if err != nil {
		return fields{}, fmt.Errorf("%s id: %w", kind, err)
	}
	weight, err := strconv.ParseFloat(parts[5], 64)
	if err != nil {
		return fields{}, fmt.Errorf("%s weight: %w", kind, err)
	}
	return fields{
		id:     id,
		name:   parts[2],
		field3: parts[3],
		field4: parts[4],
		weight: weight,
	}, nil
}

// parseIntField разбирает целочисленное поле с понятным сообщением.
func parseIntField(value, kind, name string) (int, error) {
	n, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("%s %s: %w", kind, name, err)
	}
	return n, nil
}

// ============================================================
// Weapon
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

func (w *Weapon) Serialize(out io.Writer) error {
	if _, err := fmt.Fprintf(out, "Weapon|%d|%s|%d|%d|%v\n",
		w.id, w.name, w.Damage, w.Durability, w.weight); err != nil {
		return fmt.Errorf("Weapon: %w", err)
	}
	return nil
}

func (w *Weapon) Deserialize(in io.Reader) error {
	line, err := readLine(in)
	if err != nil && !isLastLine(line, err) {
		return fmt.Errorf("Weapon: чтение строки: %w", err)
	}
	f, err := parseFields(line, "Weapon")
	if err != nil {
		return err
	}
	damage, err := parseIntField(f.field3, "Weapon", "Damage")
	if err != nil {
		return err
	}
	dur, err := parseIntField(f.field4, "Weapon", "Durability")
	if err != nil {
		return err
	}
	w.id = f.id
	w.name = f.name
	w.Damage = damage
	w.Durability = dur
	w.weight = f.weight
	return nil
}

// ============================================================
// Armor
// ============================================================

type Armor struct {
	id         int
	name       string
	weight     float64
	Defense    int
	Durability int
}

func (a *Armor) Name() string    { return a.name }
func (a *Armor) Weight() float64 { return a.weight }
func (a *Armor) ID() int         { return a.id }
func (a *Armor) SetID(id int)    { a.id = id }

func (a *Armor) Use() (string, error) {
	if a.Durability <= 0 {
		return "", fmt.Errorf("броня %q сломана: %w", a.name, ErrBroken)
	}
	a.Durability--
	return fmt.Sprintf("Надел %s (защита %d, прочность %d)", a.name, a.Defense, a.Durability), nil
}

func (a *Armor) Serialize(out io.Writer) error {
	if _, err := fmt.Fprintf(out, "Armor|%d|%s|%d|%d|%v\n",
		a.id, a.name, a.Defense, a.Durability, a.weight); err != nil {
		return fmt.Errorf("Armor: %w", err)
	}
	return nil
}

func (a *Armor) Deserialize(in io.Reader) error {
	line, err := readLine(in)
	if err != nil && !isLastLine(line, err) {
		return fmt.Errorf("Armor: чтение строки: %w", err)
	}
	f, err := parseFields(line, "Armor")
	if err != nil {
		return err
	}
	defense, err := parseIntField(f.field3, "Armor", "Defense")
	if err != nil {
		return err
	}
	dur, err := parseIntField(f.field4, "Armor", "Durability")
	if err != nil {
		return err
	}
	a.id = f.id
	a.name = f.name
	a.Defense = defense
	a.Durability = dur
	a.weight = f.weight
	return nil
}

// ============================================================
// Potion
// ============================================================

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

func (p *Potion) Serialize(out io.Writer) error {
	if _, err := fmt.Fprintf(out, "Potion|%d|%s|%s|%d|%v\n",
		p.id, p.name, p.Effect, p.Charges, p.weight); err != nil {
		return fmt.Errorf("Potion: %w", err)
	}
	return nil
}

func (p *Potion) Deserialize(in io.Reader) error {
	line, err := readLine(in)
	if err != nil && !isLastLine(line, err) {
		return fmt.Errorf("Potion: чтение строки: %w", err)
	}
	f, err := parseFields(line, "Potion")
	if err != nil {
		return err
	}
	charges, err := parseIntField(f.field4, "Potion", "Charges")
	if err != nil {
		return err
	}
	p.id = f.id
	p.name = f.name
	p.Effect = f.field3
	p.Charges = charges
	p.weight = f.weight
	return nil
}

// ============================================================
// Конструкторы
// ============================================================

// NewWeapon создаёт оружие.
func NewWeapon(name string, weight float64, damage, durability int) *Weapon {
	return &Weapon{name: name, weight: weight, Damage: damage, Durability: durability}
}

// NewArmor создаёт броню.
func NewArmor(name string, weight float64, defense, durability int) *Armor {
	return &Armor{name: name, weight: weight, Defense: defense, Durability: durability}
}

// NewPotion создаёт зелье.
func NewPotion(name string, weight float64, effect string, charges int) *Potion {
	return &Potion{name: name, weight: weight, Effect: effect, Charges: charges}
}

// ============================================================
// Фабрика по типу записи
// ============================================================

// NewByKind создаёт пустой предмет нужного типа по строковому
// идентификатору из файла ("Weapon"/"Armor"/"Potion").
func NewByKind(kind string) (Item, error) {
	switch kind {
	case "Weapon":
		return &Weapon{}, nil
	case "Armor":
		return &Armor{}, nil
	case "Potion":
		return &Potion{}, nil
	default:
		return nil, fmt.Errorf("неизвестный тип: %q", kind)
	}
}

// KindName возвращает короткое имя типа предмета ("Weapon"/"Armor"/
// "Potion") — то же, что раньше давал %T в монолитном main.
func KindName(it Item) string {
	switch it.(type) {
	case *Weapon:
		return "Weapon"
	case *Armor:
		return "Armor"
	case *Potion:
		return "Potion"
	default:
		return fmt.Sprintf("%T", it)
	}
}

// ============================================================
// Чтение строки
// ============================================================

// readLine читает ровно одну строку из r и возвращает её вместе с ошибкой.
//
// Предмет НЕ хранит ридер внутри себя: r приходит параметром в каждый вызов
// Deserialize. Поэтому один и тот же предмет можно безопасно десериализовать
// из разных источников — второй вызов прочитает из нового ридера, а не из
// «запомненного» первого.
//
// Читаем по одному байту до '\n', не полагаясь на конкретный тип ридера
// (никаких type-assertion на *bufio.Reader / *bytes.Buffer). Так мы не
// забираем из потока лишние байты следующей записи: если на вход пришёл
// общий поток, его владелец (например, хранилище) сам решает, как
// буферизовать чтение.
func readLine(r io.Reader) (string, error) {
	var sb strings.Builder
	var buf [1]byte
	for {
		n, err := r.Read(buf[:])
		if n == 1 {
			if buf[0] == '\n' {
				// Строка прочитана целиком. Если Read вернул '\n'
				// вместе с ошибкой — не теряем её: возвращаем строку
				// и эту ошибку. Вызывающий (через isLastLine) сам решит,
				// что это последняя строка без завершающего '\n'.
				return sb.String(), err
			}
			sb.WriteByte(buf[0])
		}
		if err != nil {
			return sb.String(), err
		}
	}
}

func isLastLine(line string, err error) bool {
	return errors.Is(err, io.EOF) && strings.TrimSpace(line) != ""
}

// Компиляторные проверки: все типы реализуют оба интерфейса.
var (
	_ Item     = (*Weapon)(nil)
	_ Item     = (*Armor)(nil)
	_ Item     = (*Potion)(nil)
	_ Storable = (*Weapon)(nil)
	_ Storable = (*Armor)(nil)
	_ Storable = (*Potion)(nil)
)
