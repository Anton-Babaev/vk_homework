package main

import "fmt"

type Profile struct {
	Name       string
	Age        int
	Email      string
	Newsletter bool
}

type ProfilePatch struct {
	Name       *string
	Age        *int
	Email      *string
	Newsletter *bool
}

func ApplyPatch(p *Profile, patch ProfilePatch) {
	if p == nil {
		return
	}
	if patch.Name != nil {
		p.Name = *patch.Name
	}
	if patch.Age != nil {
		p.Age = *patch.Age
	}
	if patch.Email != nil {
		p.Email = *patch.Email
	}
	if patch.Newsletter != nil {
		p.Newsletter = *patch.Newsletter
	}
}

func main() {
	profile := Profile{
		Name:       "Анна",
		Age:        30,
		Email:      "anna@example.com",
		Newsletter: true,
	}
	fmt.Println("до:   ", profile)

	newAge := 31
	newNewsletter := false
	ApplyPatch(&profile, ProfilePatch{
		Age:        &newAge,
		Newsletter: &newNewsletter,
	})
	fmt.Println("после:", profile)

	ApplyPatch(&profile, ProfilePatch{})
	fmt.Println("пустой патч:", profile)

	ApplyPatch(nil, ProfilePatch{Age: &newAge})
	fmt.Println("ApplyPatch(nil, ...) не упал")

	fmt.Println("Вопрос: почему в патче нельзя просто объявить Newsletter bool?")
	fmt.Println("Ответ: потосу что у bool всего 2 состояния, а нам надо 3")
}
