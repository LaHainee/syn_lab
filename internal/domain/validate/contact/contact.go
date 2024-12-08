package contact

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"contacts/internal/model"
)

var errValidation = errors.New("validation error")

type Validator struct{}

func New() *Validator {
	return &Validator{}
}

// Validate – валидирует поля модели Contact.
//
// Возвращает мапу полей с соответствующей ошибкой.
func (v *Validator) Validate(contact model.ContactForCreate) map[model.Field]string {
	fieldMsgs := make(map[model.Field]string)

	msg, err := v.name(contact.Name)
	if errors.Is(err, errValidation) {
		fieldMsgs[model.FieldName] = msg
	}

	msg, err = v.surname(contact.Surname)
	if errors.Is(err, errValidation) {
		fieldMsgs[model.FieldSurname] = msg
	}

	msg, err = v.birthday(contact.Birthday)
	if errors.Is(err, errValidation) {
		fieldMsgs[model.FieldBirthday] = msg
	}

	msg, err = v.phone(contact.Phone)
	if errors.Is(err, errValidation) {
		fieldMsgs[model.FieldPhone] = msg
	}

	msg, err = v.email(contact.Email)
	if errors.Is(err, errValidation) {
		fieldMsgs[model.FieldEmail] = msg
	}

	for link, value := range contact.Links {
		msg, err = v.link(link, value)
		if errors.Is(err, errValidation) {
			fieldMsgs[model.Field(link)] = msg
		}
	}

	return fieldMsgs
}

func (v *Validator) link(link model.ContactLink, value string) (string, error) {
	re := regexp.MustCompile(`^(https?://[a-zA-Z0-9.-]+(?:/[^\s]*)?)$`)

	ok := re.MatchString(value)
	if !ok {
		return fmt.Sprintf("Ссылка %s некорректная.\nФормат: https://ya.ru", string(link)), errValidation
	}

	return "", nil
}

func (v *Validator) email(email string) (string, error) {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	ok := re.MatchString(email)
	if !ok {
		return "Некорректный email", errValidation
	}

	return "", nil
}

func (v *Validator) phone(phone string) (string, error) {
	_, err := model.NewPhone(phone)
	if err != nil {
		return "Телефон должен иметь формат +7 (915) 159-67-81", errValidation
	}

	return "", nil
}

func (v *Validator) name(name string) (string, error) {
	re := regexp.MustCompile(`^[А-Яа-яЁё]{2,10}$`)

	ok := re.MatchString(name)
	if !ok {
		return "Имя должно состоять только из русских букв\nи иметь длину от 2 до 10 символов", errValidation
	}

	return "", nil
}

func (v *Validator) surname(surname string) (string, error) {
	re := regexp.MustCompile(`^[А-Яа-яЁё]{2,10}$`)

	ok := re.MatchString(surname)
	if !ok {
		return "Фамилия должна состоять только из русских букв\nи иметь длину от 2 до 10 символов", errValidation
	}

	return "", nil
}

func (v *Validator) birthday(birthday string) (string, error) {
	t, err := time.Parse("02.01.2006", birthday)
	if err != nil {
		return "Дата рождения должна быть в формате 10.01.2001", errValidation
	}

	if t.After(time.Now()) {
		return "Дата рождения не может быть в будущем", errValidation
	}

	// Минимальная дата рождения
	minBirthday := time.Date(1925, 1, 1, 0, 0, 0, 0, time.UTC)
	if t.Before(minBirthday) {
		return fmt.Sprintf("Минимальная дата рождения %s", minBirthday.Format("02.01.2006")), errValidation
	}

	return "", nil
}
