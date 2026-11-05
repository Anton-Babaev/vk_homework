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
