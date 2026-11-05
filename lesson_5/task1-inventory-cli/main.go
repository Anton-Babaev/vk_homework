package main

import (
	"errors"
	"fmt"
)

var (
	ErrBroken    = errors.New("Item is broken")
	ErrExhausted = errors.New("Item is exhausted")
)

type Weapon struct {
	ID         int
	name       string
	weight     float64
	Damage     int
	Durability int
}

type Armor struct {
	ID      int
	name    string
	weight  float64
	Defense int
}

type Potion struct {
	ID      int
	name    string
	weight  float64
	Effect  string
	Charges int
}

type Item interface {
	Name() string
	Weight() float64
	Use() (string, error)
}

func (w *Weapon) Name() string    { return w.name }
func (w *Weapon) Weight() float64 { return w.weight }
func (w *Weapon) Use() (string, error) {
	if w.Durability <= 0 {
		return "", ErrBroken
	}
	w.Durability--
	return fmt.Sprintf("В левой руке %s, нанес урона %d, осталось прочности %d", w.name, w.Damage,
		w.Durability), nil
}

func (a *Armor) Name() string    { return a.name }
func (a *Armor) Weight() float64 { return a.weight }
func (a *Armor) Use() (string, error) {
	return fmt.Sprintf("В правой руке держу %s, защита %d", a.name, a.Defense), nil
}

func (p *Potion) Name() string    { return p.name }
func (p *Potion) Weight() float64 { return p.weight }
func (p *Potion) Use() (string, error) {
	if p.Charges <= 0 {
		return "", ErrExhausted
	}
	p.Charges--
	return fmt.Sprintf("С размаху влупил %s (%s, осталось зарядов: %d)", p.name, p.Effect, p.Charges), nil
}

var (
	_ Item = (*Weapon)(nil)
	_ Item = (*Armor)(nil)
	_ Item = (*Potion)(nil)
)

func isNil(i Item) bool {
	if i == nil {
		return true
	}
	switch v := i.(type) {
	case *Weapon:
		return v == nil
	case *Armor:
		return v == nil
	case *Potion:
		return v == nil
	}
	return false
}

func DescribeItem(i Item) string {
	if isNil(i) {
		return "Предмет отсутствует"
	}
	return fmt.Sprintf("%s (вес: %.2f)", i.Name(), i.Weight())
}

func Filter[T any](items []T, predicate func(T) bool) []T {
	// Не до конца разобрался с сигнатурой функции
	result := make([]T, 0, len(items))
	for _, item := range items {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

func Map[T any, R any](items []T, transform func(T) R) []R {
	// Не до конца разобрался с сигнатурой
	result := make([]R, 0, len(items))
	for _, item := range items {
		result = append(result, transform(item))
	}
	return result
}

func Find[T any](items []T, condition func(T) bool) (T, bool) {
	// Не до конца разобрался с сигнатурой
	for _, item := range items {
		if condition(item) {
			return item, true
		}
	}
	var zero T
	return zero, false
}

type Inventory struct {
	Items  []Item
	nextID int
}

func NewInventory() *Inventory {
	return &Inventory{nextID: 1}
}

func (inv *Inventory) AddItem(item Item) int {
	id := inv.nextID
	inv.nextID++
	switch v := item.(type) {
	case *Weapon:
		v.ID = id
	case *Armor:
		v.ID = id
	case *Potion:
		v.ID = id
	}
	inv.Items = append(inv.Items, item)
	return id
}

func (inv *Inventory) FindByID(id int) (Item, bool) {
	for _, item := range inv.Items {
		var itemID int
		switch v := item.(type) {
		case *Weapon:
			itemID = v.ID
		case *Armor:
			itemID = v.ID
		case *Potion:
			itemID = v.ID
		}
		if itemID == id {
			return item, true
		}
	}
	return nil, false
}

func (inv *Inventory) DeleteByID(id int) bool {
	for i, item := range inv.Items {
		var itemID int
		switch v := item.(type) {
		case *Weapon:
			itemID = v.ID
		case *Armor:
			itemID = v.ID
		case *Potion:
			itemID = v.ID
		}
		if itemID == id {
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
	return Map(filtered, func(i Item) *Weapon {
		return i.(*Weapon)
	})
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
	return Map(inv.Items, func(i Item) string {
		return i.Name()
	})
}

func (inv *Inventory) FindItemByName(name string) (Item, bool) {
	return Find(inv.Items, func(i Item) bool {
		return i.Name() == name
	})
}
