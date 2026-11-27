package item

import (
	"errors"
	"strings"
	"testing"
)

func TestSerializeDeserializeRoundTrip(t *testing.T) {
	t.Run("weapon", func(t *testing.T) {
		src := NewWeapon("Меч", 2.5, 10, 3)
		var buf strings.Builder
		if err := src.Serialize(&buf); err != nil {
			t.Fatalf("Serialize: %v", err)
		}
		dst := &Weapon{}
		if err := dst.Deserialize(strings.NewReader(buf.String())); err != nil {
			t.Fatalf("Deserialize: %v", err)
		}
		if dst.Name() != src.Name() || dst.Weight() != src.Weight() ||
			dst.Damage != src.Damage || dst.Durability != src.Durability {
			t.Fatalf("round-trip не совпал:\n got  %+v\n want %+v", dst, src)
		}
	})

	t.Run("armor", func(t *testing.T) {
		src := NewArmor("Щит", 5, 8, 2)
		var buf strings.Builder
		if err := src.Serialize(&buf); err != nil {
			t.Fatalf("Serialize: %v", err)
		}
		dst := &Armor{}
		if err := dst.Deserialize(strings.NewReader(buf.String())); err != nil {
			t.Fatalf("Deserialize: %v", err)
		}
		if dst.Name() != src.Name() || dst.Weight() != src.Weight() ||
			dst.Defense != src.Defense || dst.Durability != src.Durability {
			t.Fatalf("round-trip не совпал:\n got  %+v\n want %+v", dst, src)
		}
	})

	t.Run("potion", func(t *testing.T) {
		src := NewPotion("Зелье", 0.5, "heal", 2)
		var buf strings.Builder
		if err := src.Serialize(&buf); err != nil {
			t.Fatalf("Serialize: %v", err)
		}
		dst := &Potion{}
		if err := dst.Deserialize(strings.NewReader(buf.String())); err != nil {
			t.Fatalf("Deserialize: %v", err)
		}
		if dst.Name() != src.Name() || dst.Weight() != src.Weight() ||
			dst.Effect != src.Effect || dst.Charges != src.Charges {
			t.Fatalf("round-trip не совпал:\n got  %+v\n want %+v", dst, src)
		}
	})
}

func TestDeserialize_LastLineWithoutNewline(t *testing.T) {
	w := &Weapon{}
	if err := w.Deserialize(strings.NewReader("Weapon|1|Меч|10|3|2.5")); err != nil {
		t.Fatalf("Deserialize без \\n: %v", err)
	}
	if w.Name() != "Меч" {
		t.Fatalf("имя = %q, ожидалось %q", w.Name(), "Меч")
	}
}

func TestDeserialize_BrokenLine(t *testing.T) {
	w := &Weapon{}
	if err := w.Deserialize(strings.NewReader("Weapon|1|Меч\n")); err == nil {
		t.Fatal("ожидалась ошибка о неверном числе полей")
	}
}

func TestUse_WeaponBreaks(t *testing.T) {
	w := NewWeapon("Меч", 2.5, 10, 1)
	if _, err := w.Use(); err != nil {
		t.Fatalf("первое использование: %v", err)
	}
	_, err := w.Use()
	if !errors.Is(err, ErrBroken) {
		t.Fatalf("ожидалась ErrBroken, получено %v", err)
	}
}

func TestUse_PotionExhausted(t *testing.T) {
	p := NewPotion("Зелье", 0.5, "heal", 1)
	if _, err := p.Use(); err != nil {
		t.Fatalf("первое использование: %v", err)
	}
	_, err := p.Use()
	if !errors.Is(err, ErrExhausted) {
		t.Fatalf("ожидалась ErrExhausted, получено %v", err)
	}
}
