package contacts_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "contacts/internal/domain/contacts"
	"contacts/internal/model"
)

func TestAllowedLinks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		expectations func(t assert.TestingT, actual []model.ContactLink)
	}{
		{
			name: "Success",
			expectations: func(t assert.TestingT, actual []model.ContactLink) {
				expected := []model.ContactLink{
					model.ContactLinkVk,
				}

				assert.Equal(t, expected, actual)
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			out := AllowedLinks()

			tc.expectations(t, out)
		})
	}
}
