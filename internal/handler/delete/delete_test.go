package delete_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	. "contacts/internal/handler/delete"
)

func TestHandler_Delete(t *testing.T) {
	t.Parallel()

	const uuid = "uuid"

	tests := []struct {
		name         string
		uuid         string
		prepare      func(storage *Mockstorage)
		expectations func(t assert.TestingT, err error)
	}{
		{
			name: "Failed to delete contact from storage",
			uuid: uuid,
			prepare: func(storage *Mockstorage) {
				storage.EXPECT().
					Delete(uuid).
					Return(assert.AnError)
			},
			expectations: func(t assert.TestingT, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "Success",
			uuid: uuid,
			prepare: func(storage *Mockstorage) {
				storage.EXPECT().
					Delete(uuid).
					Return(nil)
			},
			expectations: func(t assert.TestingT, err error) {
				assert.NoError(t, err)
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

			err := instance.Delete(context.Background(), tc.uuid)

			tc.expectations(t, err)
		})
	}
}
