package inventory

import (
	"errors"
	"strings"
	"testing"

	"github.com/Anton-Babaev/vk_homework/lesson_6/task1-inventory-cli/internal/item"
)

// TestLoad_FileEndings проверяет границы конца файла — именно здесь
// раньше был баг: файл, заканчивающийся на '\n', ошибочно считался
// повреждённым.
func TestLoad_FileEndings(t *testing.T) {
	const weapon = "Weapon|1|Меч|10|3|2.5"

	cases := []struct {
		name  string
		input string
		want  int
	}{
		{"с завершающим \\n", weapon + "\n", 1},
		{"без завершающего \\n", weapon, 1},
		{"две строки с \\n", weapon + "\nArmor|2|Щит|8|2|5\n", 2},
		{"две строки без \\n в конце", weapon + "\nArmor|2|Щит|8|2|5", 2},
		{"пустой ввод", "", 0},
		{"только \\n", "\n", 0},
		{"пустые строки между записями", weapon + "\n\n\nArmor|2|Щит|8|2|5\n", 2},
		{"CRLF", weapon + "\r\nArmor|2|Щит|8|2|5\r\n", 2},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			inv := New()
			if err := inv.Load(strings.NewReader(tc.input)); err != nil {
				t.Fatalf("Load вернул ошибку: %v", err)
			}
			if inv.Len() != tc.want {
				t.Fatalf("загружено %d предметов, ожидалось %d", inv.Len(), tc.want)
			}
		})
	}
}

// TestLoad_BrokenLineKeepsValid проверяет политику «битая строка не
// прерывает загрузку»: валидные записи остаются, ошибка возвращается.
func TestLoad_BrokenLineKeepsValid(t *testing.T) {
	input := "Weapon|1|Меч|10|3|2.5\n" +
		"garbage line\n" +
		"Armor|2|Щит|8|2|5\n"

	inv := New()
	err := inv.Load(strings.NewReader(input))
	if err == nil {
		t.Fatal("ожидалась ошибка о битой строке")
	}
	var pie *ParseItemError
	if !errors.As(err, &pie) {
		t.Fatalf("ожидалась ParseItemError, получено %T", err)
	}
	if pie.Line != 2 {
		t.Fatalf("номер битой строки = %d, ожидалось 2", pie.Line)
	}
	if inv.Len() != 2 {
		t.Fatalf("валидных предметов %d, ожидалось 2", inv.Len())
	}
}

// TestLoad_DuplicateID проверяет, что дубликат ID — ошибка и не
// добавляется повторно.
func TestLoad_DuplicateID(t *testing.T) {
	input := "Weapon|1|Меч|10|3|2.5\n" +
		"Weapon|1|Другой меч|5|1|1.0\n"

	inv := New()
	err := inv.Load(strings.NewReader(input))
	if err == nil {
		t.Fatal("ожидалась ошибка о дубликате ID")
	}
	if inv.Len() != 1 {
		t.Fatalf("предметов %d, ожидался 1 (дубликат пропущен)", inv.Len())
	}
}

// TestLoad_RepeatedLoadNoDuplicates проверяет, что повторный Load на
// том же инвентаре не создаёт дубликаты (seen инициализируется
// уже загруженными ID).
func TestLoad_RepeatedLoadNoDuplicates(t *testing.T) {
	const input = "Weapon|1|Меч|10|3|2.5\n"

	inv := New()
	if err := inv.Load(strings.NewReader(input)); err != nil {
		t.Fatalf("первый Load: %v", err)
	}
	if err := inv.Load(strings.NewReader(input)); err == nil {
		t.Fatal("второй Load должен вернуть ошибку о дубликате")
	}
	if inv.Len() != 1 {
		t.Fatalf("предметов %d, ожидался 1", inv.Len())
	}
}

// TestSaveLoadRoundTrip проверяет, что Save→Load сохраняет данные.
func TestSaveLoadRoundTrip(t *testing.T) {
	src := New()
	src.AddItem(item.NewWeapon("Меч", 2.5, 10, 3))
	src.AddItem(item.NewArmor("Щит", 5, 8, 2))
	src.AddItem(item.NewPotion("Зелье", 0.5, "heal", 2))

	var buf strings.Builder
	if err := src.Save(&buf); err != nil {
		t.Fatalf("Save: %v", err)
	}

	dst := New()
	if err := dst.Load(strings.NewReader(buf.String())); err != nil {
		t.Fatalf("Load: %v", err)
	}
	if dst.Len() != src.Len() {
		t.Fatalf("после round-trip %d предметов, ожидалось %d", dst.Len(), src.Len())
	}
	for i, it := range src.Items() {
		assertItemEqual(t, it, dst.Items()[i])
	}
}

// assertItemEqual проверяет, что после Save→Load совпали и базовые поля
// (ID/Name/Weight), и специфичные (Damage/Durability/Defense/Effect/Charges).
// Без последних регрессия в разборе field3/field4 осталась бы незамеченной.
func assertItemEqual(t *testing.T, want, got item.Item) {
	t.Helper()
	if got.ID() != want.ID() || got.Name() != want.Name() || got.Weight() != want.Weight() {
		t.Fatalf("базовые поля не совпали: got %+v, want %+v", got, want)
	}
	switch w := want.(type) {
	case *item.Weapon:
		g, ok := got.(*item.Weapon)
		if !ok {
			t.Fatalf("тип не совпал: got %T, want %T", got, want)
		}
		if g.Damage != w.Damage || g.Durability != w.Durability {
			t.Fatalf("поля Weapon не совпали: got %+v, want %+v", g, w)
		}
	case *item.Armor:
		g, ok := got.(*item.Armor)
		if !ok {
			t.Fatalf("тип не совпал: got %T, want %T", got, want)
		}
		if g.Defense != w.Defense || g.Durability != w.Durability {
			t.Fatalf("поля Armor не совпали: got %+v, want %+v", g, w)
		}
	case *item.Potion:
		g, ok := got.(*item.Potion)
		if !ok {
			t.Fatalf("тип не совпал: got %T, want %T", got, want)
		}
		if g.Effect != w.Effect || g.Charges != w.Charges {
			t.Fatalf("поля Potion не совпали: got %+v, want %+v", g, w)
		}
	}
}

// TestAddItem_NextIDAfterLoad проверяет, что после загрузки новые ID
// продолжают нумерацию, а не начинаются с 1.
func TestAddItem_NextIDAfterLoad(t *testing.T) {
	inv := New()
	if err := inv.Load(strings.NewReader("Weapon|5|Меч|10|3|2.5\n")); err != nil {
		t.Fatalf("Load: %v", err)
	}
	id := inv.AddItem(item.NewArmor("Щит", 5, 8, 2))
	if id != 6 {
		t.Fatalf("новый ID = %d, ожидался 6", id)
	}
}
