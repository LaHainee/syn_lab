package fetch_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	. "contacts/internal/handler/fetch"
	"contacts/internal/model"
)

func TestHandler_Fetch(t *testing.T) {
	t.Parallel()

	contact := model.Contact{
		UUID:     "uuid",
		Surname:  "Ершов",
		Name:     "Виталий",
		Birthday: time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
		Phone:    model.NewPhoneFromInt64(79151596781),
		Email:    "vaershov@avito.ru",
		Links: map[model.ContactLink]string{
			model.ContactLinkVk: "vk.com",
		},
	}

	tests := []struct {
		name         string
		prepare      func(storage *Mockstorage)
		expectations func(t assert.TestingT, actual []model.Contact, err error)
	}{
		{
			name: "Failed to fetch",
			prepare: func(storage *Mockstorage) {
				storage.EXPECT().
					Fetch().
					Return(nil, assert.AnError)
			},
			expectations: func(t assert.TestingT, actual []model.Contact, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "Success",
			prepare: func(storage *Mockstorage) {
				storage.EXPECT().
					Fetch().
					Return([]model.Contact{
						contact,
					}, nil)
			},
			expectations: func(t assert.TestingT, actual []model.Contact, err error) {
				assert.NoError(t, err)
				assert.Equal(t, []model.Contact{contact}, actual)
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mockStorage := NewMockstorage(ctrl)

			if tc.prepare != nil {
				tc.prepare(mockStorage)
			}

			instance := NewHandler(mockStorage)

			out, err := instance.Fetch(context.Background())

			tc.expectations(t, out, err)
		})
	}
}

func TestHandler_FetchByUuid(t *testing.T) {
	t.Parallel()

	const uuid = "uuid"

	contact := model.Contact{
		UUID:     "uuid",
		Surname:  "Ершов",
		Name:     "Виталий",
		Birthday: time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
		Phone:    model.NewPhoneFromInt64(79151596781),
		Email:    "vaershov@avito.ru",
		Links: map[model.ContactLink]string{
			model.ContactLinkVk: "vk.com",
		},
	}

	tests := []struct {
		name         string
		uuid         string
		prepare      func(storage *Mockstorage)
		expectations func(t assert.TestingT, actual model.Contact, err error)
	}{
		{
			name: "Failed to fetch",
			uuid: uuid,
			prepare: func(storage *Mockstorage) {
				storage.EXPECT().
					FetchByUuid(uuid).
					Return(model.Contact{}, assert.AnError)
			},
			expectations: func(t assert.TestingT, actual model.Contact, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "Success",
			uuid: uuid,
			prepare: func(storage *Mockstorage) {
				storage.EXPECT().
					FetchByUuid(uuid).
					Return(contact, nil)
			},
			expectations: func(t assert.TestingT, actual model.Contact, err error) {
				assert.NoError(t, err)
				assert.Equal(t, contact, actual)
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mockStorage := NewMockstorage(ctrl)

			if tc.prepare != nil {
				tc.prepare(mockStorage)
			}

			instance := NewHandler(mockStorage)

			out, err := instance.FetchByUuid(context.Background(), tc.uuid)

			tc.expectations(t, out, err)
		})
	}
}
