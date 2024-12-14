package contact_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	. "contacts/internal/domain/validate/contact"
	"contacts/internal/model"
	"contacts/util/pointer"
)

func TestValidator_Validate(t *testing.T) {
	t.Parallel()

	valid := model.ContactForCreate{
		UUID:     pointer.To("1"),
		Name:     "Виталий",
		Surname:  "Ершов",
		Birthday: "10.01.2001",
		Phone:    "+7 (915) 159-67-81",
		Email:    "vaershov@avito.ru",
		Links: map[model.ContactLink]string{
			model.ContactLinkVk: "https://vk.com",
		},
	}

	tests := []struct {
		name         string
		contact      func() model.ContactForCreate
		expectations func(t assert.TestingT, actual map[model.Field]string)
	}{
		{
			name: "Invalid name",
			contact: func() model.ContactForCreate {
				// Перезапишем только нужное поле
				copied := valid
				copied.Name = "Hello world from vaershov"
				return copied
			},
			expectations: func(t assert.TestingT, actual map[model.Field]string) {
				expected := map[model.Field]string{
					model.FieldName: "Имя должно состоять только из русских букв\nи иметь длину от 2 до 10 символов",
				}

				assert.Equal(t, expected, actual)
			},
		},
		{
			name: "Invalid surname",
			contact: func() model.ContactForCreate {
				// Перезапишем только нужное поле
				copied := valid
				copied.Surname = "Hello world from vaershov"
				return copied
			},
			expectations: func(t assert.TestingT, actual map[model.Field]string) {
				expected := map[model.Field]string{
					model.FieldSurname: "Фамилия должна состоять только из русских букв\nи иметь длину от 2 до 10 символов",
				}

				assert.Equal(t, expected, actual)
			},
		},
		{
			name: "Invalid birthday, failed to parse as 02.01.2006",
			contact: func() model.ContactForCreate {
				// Перезапишем только нужное поле
				copied := valid
				copied.Birthday = "2001-01-01"
				return copied
			},
			expectations: func(t assert.TestingT, actual map[model.Field]string) {
				expected := map[model.Field]string{
					model.FieldBirthday: "Дата рождения должна быть в формате 10.01.2001",
				}

				assert.Equal(t, expected, actual)
			},
		},
		{
			name: "Invalid birthday, in future",
			contact: func() model.ContactForCreate {
				now := time.Now()
				now = now.Add(24 * time.Hour)

				// Перезапишем только нужное поле
				copied := valid
				copied.Birthday = now.Format("02.01.2006")
				return copied
			},
			expectations: func(t assert.TestingT, actual map[model.Field]string) {
				expected := map[model.Field]string{
					model.FieldBirthday: "Дата рождения не может быть в будущем",
				}

				assert.Equal(t, expected, actual)
			},
		},
		{
			name: "Invalid birthday, before min date",
			contact: func() model.ContactForCreate {
				// Перезапишем только нужное поле
				copied := valid
				copied.Birthday = "10.01.1910"
				return copied
			},
			expectations: func(t assert.TestingT, actual map[model.Field]string) {
				expected := map[model.Field]string{
					model.FieldBirthday: "Минимальная дата рождения 01.01.1925",
				}

				assert.Equal(t, expected, actual)
			},
		},
		{
			name: "Invalid phone",
			contact: func() model.ContactForCreate {
				// Перезапишем только нужное поле
				copied := valid
				copied.Phone = "1234"
				return copied
			},
			expectations: func(t assert.TestingT, actual map[model.Field]string) {
				expected := map[model.Field]string{
					model.FieldPhone: "Телефон должен иметь формат +7 (915) 159-67-81",
				}

				assert.Equal(t, expected, actual)
			},
		},
		{
			name: "Invalid email",
			contact: func() model.ContactForCreate {
				// Перезапишем только нужное поле
				copied := valid
				copied.Email = "1234"
				return copied
			},
			expectations: func(t assert.TestingT, actual map[model.Field]string) {
				expected := map[model.Field]string{
					model.FieldEmail: "Некорректный email",
				}

				assert.Equal(t, expected, actual)
			},
		},
		{
			name: "Invalid link",
			contact: func() model.ContactForCreate {
				// Перезапишем только нужное поле
				copied := valid
				copied.Links = map[model.ContactLink]string{
					model.ContactLinkVk: "hello",
				}
				return copied
			},
			expectations: func(t assert.TestingT, actual map[model.Field]string) {
				expected := map[model.Field]string{
					model.Field(model.ContactLinkVk): "Ссылка vk.com некорректная.\nФормат: https://ya.ru",
				}

				assert.Equal(t, expected, actual)
			},
		},
		{
			name: "Valid contact",
			contact: func() model.ContactForCreate {
				return valid
			},
			expectations: func(t assert.TestingT, actual map[model.Field]string) {
				assert.Empty(t, actual)
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			instance := New()

			out := instance.Validate(tc.contact())

			tc.expectations(t, out)
		})
	}
}
