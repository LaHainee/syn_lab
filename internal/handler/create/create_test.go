package create_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	. "contacts/internal/handler/create"
	"contacts/internal/model"
)

func TestHandler_Create(t *testing.T) {
	t.Parallel()

	contact := model.ContactForCreate{
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
		prepare          func(storage *Mockstorage, validator *Mockvalidator, uuid *Mockuuid)
		expectations     func(t assert.TestingT, actual map[model.Field]string, err error)
	}{
		{
			name:             "Validation error",
			contactForCreate: contact,
			prepare: func(_ *Mockstorage, validator *Mockvalidator, _ *Mockuuid) {
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
				Birthday: "2001-01-10",
				Phone:    "+7 (915) 159-67-81",
			},
			prepare: func(_ *Mockstorage, validator *Mockvalidator, _ *Mockuuid) {
				validator.EXPECT().
					Validate(model.ContactForCreate{
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
				Birthday: "10.01.2001",
				Phone:    "+7 (915) 1596781",
			},
			prepare: func(_ *Mockstorage, validator *Mockvalidator, _ *Mockuuid) {
				validator.EXPECT().
					Validate(model.ContactForCreate{
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
			name:             "Failed to create in storage",
			contactForCreate: contact,
			prepare: func(storage *Mockstorage, validator *Mockvalidator, uuid *Mockuuid) {
				validator.EXPECT().
					Validate(contact).
					Return(nil)

				uuid.EXPECT().
					NewString().
					Return("uuid")

				b, err := time.Parse("02.01.2006", contact.Birthday)
				assert.NoError(t, err)

				storage.EXPECT().
					Create(model.Contact{
						UUID:     "uuid",
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
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mockStorage := NewMockstorage(ctrl)
			mockValidator := NewMockvalidator(ctrl)
			mockUuid := NewMockuuid(ctrl)

			if tc.prepare != nil {
				tc.prepare(mockStorage, mockValidator, mockUuid)
			}

			instance := NewHandler(mockStorage, mockUuid, mockValidator)

			out, err := instance.Create(context.Background(), tc.contactForCreate)

			tc.expectations(t, out, err)
		})
	}
}
