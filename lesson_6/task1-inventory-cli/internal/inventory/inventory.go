// Package inventory содержит правила работы с набором предметов:
// добавление, поиск, удаление, а также загрузку/сохранение через
// абстрактное хранилище. Пакет зависит только от item.
package inventory

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/Anton-Babaev/vk_homework/lesson_6/task1-inventory-cli/internal/item"
)

// ParseItemError — типизированная ошибка разбора одной строки.
type ParseItemError struct {
	Line  int
	Value string
	Err   error
}

func (e *ParseItemError) Error() string {
	return fmt.Sprintf("строка %d: не удалось разобрать %q: %v", e.Line, e.Value, e.Err)
}

func (e *ParseItemError) Unwrap() error { return e.Err }

// Inventory — набор предметов с уникальными ID.
type Inventory struct {
	items  []item.Item
	nextID int
}

// New создаёт пустой инвентарь.
func New() *Inventory {
	return &Inventory{}
}

// Items возвращает копию среза предметов, чтобы вызывающий не мог
// менять внутреннее состояние инвентаря в обход его методов.
func (inv *Inventory) Items() []item.Item {
	out := make([]item.Item, len(inv.items))
	copy(out, inv.items)
	return out
}

// Len возвращает количество предметов.
func (inv *Inventory) Len() int { return len(inv.items) }

// AddItem присваивает предмету новый ID и добавляет его в инвентарь.
func (inv *Inventory) AddItem(it item.Item) int {
	inv.nextID++
	it.SetID(inv.nextID)
	inv.items = append(inv.items, it)
	return inv.nextID
}

// FindByID ищет предмет по ID.
func (inv *Inventory) FindByID(id int) (item.Item, bool) {
	for _, it := range inv.items {
		if it.ID() == id {
			return it, true
		}
	}
	return nil, false
}

// DeleteByID удаляет предмет по ID. Возвращает false, если не найден.
func (inv *Inventory) DeleteByID(id int) bool {
	for i, it := range inv.items {
		if it.ID() == id {
			inv.items = append(inv.items[:i], inv.items[i+1:]...)
			return true
		}
	}
	return false
}

// ============================================================
// Загрузка / сохранение
// ============================================================

// Load читает инвентарь построчно из r.
//
// Политика загрузки: повреждённые строки НЕ прерывают загрузку — они
// пропускаются, а все ошибки собираются и возвращаются объединёнными
// через errors.Join. Так пользователь видит сразу все проблемы файла,
// а валидные записи остаются доступны. Дубликаты ID также считаются
// ошибкой и пропускаются (иначе FindByID/DeleteByID работали бы
// непредсказуемо).
//
// seen инициализируется уже загруженными ID, поэтому повторный Load
// на том же инвентаре не создаёт дубликаты.
func (inv *Inventory) Load(r io.Reader) error {
	br := bufio.NewReader(r)
	var errs []error
	seen := make(map[int]bool, len(inv.items))
	for _, it := range inv.items {
		seen[it.ID()] = true
	}
	lineNo := 0

	for {
		line, err := br.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				if strings.TrimSpace(line) != "" {
					lineNo++
					inv.parseLine(line, lineNo, seen, &errs)
				}
				break
			}
			errs = append(errs, fmt.Errorf("чтение файла: %w", err))
			break
		}
		lineNo++
		line = strings.TrimRight(line, "\r\n")
		if strings.TrimSpace(line) == "" {
			continue
		}
		inv.parseLine(line, lineNo, seen, &errs)
	}

	return errors.Join(errs...)
}

func (inv *Inventory) parseLine(line string, lineNo int, seen map[int]bool, errs *[]error) {
	it, err := parseItemLine(line)
	if err != nil {
		*errs = append(*errs, &ParseItemError{Line: lineNo, Value: line, Err: err})
		return
	}
	if seen[it.ID()] {
		*errs = append(*errs, &ParseItemError{
			Line:  lineNo,
			Value: line,
			Err:   fmt.Errorf("дубликат ID=%d", it.ID()),
		})
		return
	}
	seen[it.ID()] = true
	inv.items = append(inv.items, it)
	if it.ID() > inv.nextID {
		inv.nextID = it.ID()
	}
}

// Save записывает все предметы в w.
func (inv *Inventory) Save(w io.Writer) error {
	bw := bufio.NewWriter(w)
	for _, it := range inv.items {
		s, ok := it.(item.Storable)
		if !ok {
			return fmt.Errorf("сериализация %q: тип %T не реализует Storable", it.Name(), it)
		}
		if err := s.Serialize(bw); err != nil {
			return fmt.Errorf("сериализация %q: %w", it.Name(), err)
		}
	}
	if err := bw.Flush(); err != nil {
		return fmt.Errorf("flush: %w", err)
	}
	return nil
}

// parseItemLine разбирает одну строку формата в предмет.
func parseItemLine(line string) (item.Item, error) {
	parts := strings.Split(line, "|")
	if len(parts) < 2 {
		return nil, fmt.Errorf("ожидалось минимум 2 поля, получено %d", len(parts))
	}
	it, err := item.NewByKind(parts[0])
	if err != nil {
		return nil, err
	}
	s, ok := it.(item.Storable)
	if !ok {
		return nil, fmt.Errorf("тип %T не реализует Storable", it)
	}
	if err := s.Deserialize(strings.NewReader(line + "\n")); err != nil {
		return nil, err
	}
	return it, nil
}
