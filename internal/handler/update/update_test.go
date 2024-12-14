package update_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	. "contacts/internal/handler/update"
	"contacts/internal/model"
	"contacts/util/pointer"
)

func TestHandler_Update(t *testing.T) {
	t.Parallel()

	contact := model.ContactForCreate{
		UUID:     pointer.To("1"),
		Name:     "Виталий",
		Surname:  "Ершов",
		Birthday: "10.01.2001",
		Phone:    "+7 (915) 159-67-81",
		Email:    "vaershov@avito.ru",
		Links: map[model.ContactLink]string{
			model.ContactLinkVk: "vk.com",
		},
	}

	tests := []struct {
		name             string
		contactForCreate model.ContactForCreate
		prepare          func(storage *Mockstorage, validator *Mockvalidator)
		expectations     func(t assert.TestingT, got map[model.Field]string, err error)
	}{
		{
			name:             "Validation error",
			contactForCreate: contact,
			prepare: func(_ *Mockstorage, validator *Mockvalidator) {
				validator.EXPECT().
					Validate(contact).
					Return(map[model.Field]string{
						model.FieldName: "msg",
					})
			},
			expectations: func(t assert.TestingT, actual map[model.Field]string, err error) {
				assert.ErrorIs(t, err, model.ErrValidation)

				expected := map[model.Field]string{
					model.FieldName: "msg",
				}

				assert.Equal(t, expected, actual)
			},
		},
		{
			name: "Failed to parse birthday from string",
			contactForCreate: model.ContactForCreate{
				UUID:     pointer.To("1"),
				Birthday: "2001-01-10",
				Phone:    "+7 (915) 159-67-81",
			},
			prepare: func(_ *Mockstorage, validator *Mockvalidator) {
				validator.EXPECT().
					Validate(model.ContactForCreate{
						UUID:     pointer.To("1"),
						Birthday: "2001-01-10",
						Phone:    "+7 (915) 159-67-81",
					}).
					Return(nil)
			},
			expectations: func(t assert.TestingT, actual map[model.Field]string, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "Failed to parse phone from string",
			contactForCreate: model.ContactForCreate{
				UUID:     pointer.To("1"),
				Birthday: "10.01.2001",
				Phone:    "+7 (915) 1596781",
			},
			prepare: func(_ *Mockstorage, validator *Mockvalidator) {
				validator.EXPECT().
					Validate(model.ContactForCreate{
						UUID:     pointer.To("1"),
						Birthday: "10.01.2001",
						Phone:    "+7 (915) 1596781",
					}).
					Return(nil)
			},
			expectations: func(t assert.TestingT, actual map[model.Field]string, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "Empty uuid",
			contactForCreate: model.ContactForCreate{
				Birthday: "10.01.2001",
				Phone:    "+7 (915) 159-67-81",
			},
			prepare: func(_ *Mockstorage, validator *Mockvalidator) {
				validator.EXPECT().
					Validate(model.ContactForCreate{
						Birthday: "10.01.2001",
						Phone:    "+7 (915) 159-67-81",
					}).
					Return(nil)
			},
			expectations: func(t assert.TestingT, actual map[model.Field]string, err error) {
				assert.Error(t, err)
			},
		},
		{
			name:             "Failed to update",
			contactForCreate: contact,
			prepare: func(storage *Mockstorage, validator *Mockvalidator) {
				validator.EXPECT().
					Validate(contact).
					Return(nil)

				b, err := time.Parse("02.01.2006", contact.Birthday)
				assert.NoError(t, err)

				storage.EXPECT().
					Update(model.Contact{
						UUID:     "1",
						Name:     "Виталий",
						Surname:  "Ершов",
						Birthday: b,
						Phone:    model.NewPhoneFromInt64(79151596781),
						Email:    "vaershov@avito.ru",
						Links: map[model.ContactLink]string{
							model.ContactLinkVk: "vk.com",
						},
					}).
					Return(assert.AnError)
			},
			expectations: func(t assert.TestingT, actual map[model.Field]string, err error) {
				assert.Error(t, err)
			},
		},
		{
			name:             "Success",
			contactForCreate: contact,
			prepare: func(storage *Mockstorage, validator *Mockvalidator) {
				validator.EXPECT().
					Validate(contact).
					Return(nil)

				b, err := time.Parse("02.01.2006", contact.Birthday)
				assert.NoError(t, err)

				storage.EXPECT().
					Update(model.Contact{
						UUID:     "1",
						Name:     "Виталий",
						Surname:  "Ершов",
						Birthday: b,
						Phone:    model.NewPhoneFromInt64(79151596781),
						Email:    "vaershov@avito.ru",
						Links: map[model.ContactLink]string{
							model.ContactLinkVk: "vk.com",
						},
					}).
					Return(nil)
			},
			expectations: func(t assert.TestingT, actual map[model.Field]string, err error) {
				assert.NoError(t, err)
				assert.Empty(t, actual)
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mockStorage := NewMockstorage(ctrl)
			mockValidator := NewMockvalidator(ctrl)

			if tc.prepare != nil {
				tc.prepare(mockStorage, mockValidator)
			}

			instance := NewHandler(mockStorage, mockValidator)

			out, err := instance.Update(context.Background(), tc.contactForCreate)

			tc.expectations(t, out, err)
		})
	}
}
