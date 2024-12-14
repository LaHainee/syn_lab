package search_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	. "contacts/internal/handler/search"
	"contacts/internal/model"
)

func TestHandler_Search(t *testing.T) {
	t.Parallel()

	req := model.SearchRequest{}

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
		request      model.SearchRequest
		prepare      func(storage *Mockstorage)
		expectations func(t assert.TestingT, got []model.Contact, err error)
	}{
		{
			name:    "Failed to search",
			request: req,
			prepare: func(storage *Mockstorage) {
				storage.EXPECT().
					Search(req).
					Return(nil, assert.AnError)
			},
			expectations: func(t assert.TestingT, got []model.Contact, err error) {
				assert.Error(t, err)
			},
		},
		{
			name:    "Success",
			request: req,
			prepare: func(storage *Mockstorage) {
				storage.EXPECT().
					Search(req).
					Return([]model.Contact{
						contact,
					}, nil)
			},
			expectations: func(t assert.TestingT, got []model.Contact, err error) {
				assert.NoError(t, err)
				assert.Equal(t, []model.Contact{contact}, got)
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

			out, err := instance.Search(context.Background(), tc.request)

			tc.expectations(t, out, err)
		})
	}
}
