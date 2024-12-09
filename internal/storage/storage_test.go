package storage_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"contacts/internal/model"
	. "contacts/internal/storage"
)

func TestStorage_Search(t *testing.T) {
	t.Parallel()

	req := model.SearchRequest{
		Query: "Ершов",
	}

	tests := []struct {
		name         string
		request      model.SearchRequest
		prepare      func(db *Mockdatabase)
		expectations func(t assert.TestingT, actual []model.Contact, err error)
	}{
		{
			name:    "Failed to read contacts from database",
			request: req,
			prepare: func(db *Mockdatabase) {
				db.EXPECT().
					Read().
					Return(nil, assert.AnError)
			},
			expectations: func(t assert.TestingT, actual []model.Contact, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "Empty query",
			request: model.SearchRequest{
				Query: "",
			},
			prepare: func(db *Mockdatabase) {
				db.EXPECT().
					Read().
					Return(map[string]Contact{
						"1": {
							UUID:    "1",
							Surname: "Ершов",
						},
						"2": {
							UUID:    "2",
							Surname: "Никандров",
						},
					}, nil)
			},
			expectations: func(t assert.TestingT, actual []model.Contact, err error) {
				assert.NoError(t, err)

				assert.Contains(t, actual, model.Contact{
					UUID:    "1",
					Surname: "Ершов",
					Links:   map[model.ContactLink]string{},
				})

				assert.Contains(t, actual, model.Contact{
					UUID:    "2",
					Surname: "Никандров",
					Links:   map[model.ContactLink]string{},
				})
			},
		},
		{
			name:    "Success",
			request: req,
			prepare: func(db *Mockdatabase) {
				db.EXPECT().
					Read().
					Return(map[string]Contact{
						"1": {
							UUID:    "1",
							Surname: "Ершов",
						},
						"2": {
							UUID:    "2",
							Surname: "Никандров",
						},
					}, nil)
			},
			expectations: func(t assert.TestingT, actual []model.Contact, err error) {
				assert.NoError(t, err)

				assert.Contains(t, actual, model.Contact{
					UUID:    "1",
					Surname: "Ершов",
					Links:   map[model.ContactLink]string{},
				})
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mockDatabase := NewMockdatabase(ctrl)

			if tc.prepare != nil {
				tc.prepare(mockDatabase)
			}

			instance := New(mockDatabase)

			out, err := instance.Search(tc.request)

			tc.expectations(t, out, err)
		})
	}
}

func TestStorage_FetchByUuid(t *testing.T) {
	t.Parallel()

	const uuid = "311eec60-144b-4558-9799-2d44f81cb68a"

	tests := []struct {
		name         string
		uuid         string
		prepare      func(db *Mockdatabase)
		expectations func(t assert.TestingT, actual model.Contact, err error)
	}{
		{
			name: "Failed to read from database",
			uuid: uuid,
			prepare: func(db *Mockdatabase) {
				db.EXPECT().
					Read().
					Return(nil, assert.AnError)
			},
			expectations: func(t assert.TestingT, actual model.Contact, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "Contact not found",
			uuid: uuid,
			prepare: func(db *Mockdatabase) {
				db.EXPECT().
					Read().
					Return(map[string]Contact{
						"1": {
							UUID:    "1",
							Surname: "Ершов",
						},
					}, nil)
			},
			expectations: func(t assert.TestingT, actual model.Contact, err error) {
				assert.ErrorIs(t, err, model.ErrNotFound)
			},
		},
		{
			name: "Success",
			uuid: uuid,
			prepare: func(db *Mockdatabase) {
				db.EXPECT().
					Read().
					Return(map[string]Contact{
						uuid: {
							UUID:    uuid,
							Surname: "Ершов",
						},
					}, nil)
			},
			expectations: func(t assert.TestingT, actual model.Contact, err error) {
				assert.NoError(t, err)

				expected := model.Contact{
					UUID:    uuid,
					Surname: "Ершов",
					Links:   map[model.ContactLink]string{},
				}

				assert.Equal(t, expected, actual)
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mockDatabase := NewMockdatabase(ctrl)

			if tc.prepare != nil {
				tc.prepare(mockDatabase)
			}

			instance := New(mockDatabase)

			out, err := instance.FetchByUuid(tc.uuid)

			tc.expectations(t, out, err)
		})
	}
}

func TestStorage_Fetch(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		prepare      func(db *Mockdatabase)
		expectations func(t assert.TestingT, actual []model.Contact, err error)
	}{
		{
			name: "Failed to read from database",
			prepare: func(db *Mockdatabase) {
				db.EXPECT().
					Read().
					Return(nil, assert.AnError)
			},
			expectations: func(t assert.TestingT, actual []model.Contact, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "Success",
			prepare: func(db *Mockdatabase) {
				db.EXPECT().
					Read().
					Return(map[string]Contact{
						"1": {
							UUID:    "1",
							Surname: "Ершов",
						},
					}, nil)
			},
			expectations: func(t assert.TestingT, actual []model.Contact, err error) {
				assert.NoError(t, err)

				expected := []model.Contact{
					{
						UUID:    "1",
						Surname: "Ершов",
						Links:   map[model.ContactLink]string{},
					},
				}

				assert.Equal(t, expected, actual)
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mockDatabase := NewMockdatabase(ctrl)

			if tc.prepare != nil {
				tc.prepare(mockDatabase)
			}

			instance := New(mockDatabase)

			out, err := instance.Fetch()

			tc.expectations(t, out, err)
		})
	}
}

func TestStorage_Delete(t *testing.T) {
	t.Parallel()

	const uuid = "b57d05a2-856c-4840-8842-af62165669cf"

	tests := []struct {
		name         string
		uuid         string
		prepare      func(db *Mockdatabase)
		expectations func(t assert.TestingT, err error)
	}{
		{
			name: "Failed to read from database",
			uuid: uuid,
			prepare: func(db *Mockdatabase) {
				db.EXPECT().
					Read().
					Return(nil, assert.AnError)
			},
			expectations: func(t assert.TestingT, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "Failed to save to database",
			uuid: uuid,
			prepare: func(db *Mockdatabase) {
				db.EXPECT().
					Read().
					Return(map[string]Contact{
						"1": {
							UUID: "1",
						},
						uuid: {
							UUID: uuid,
						},
					}, nil)

				db.EXPECT().
					Save(map[string]Contact{
						"1": {
							UUID: "1",
						},
					}).
					Return(assert.AnError)
			},
			expectations: func(t assert.TestingT, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "Success",
			uuid: uuid,
			prepare: func(db *Mockdatabase) {
				db.EXPECT().
					Read().
					Return(map[string]Contact{
						"1": {
							UUID: "1",
						},
						uuid: {
							UUID: uuid,
						},
					}, nil)

				db.EXPECT().
					Save(map[string]Contact{
						"1": {
							UUID: "1",
						},
					}).
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
			mockDatabase := NewMockdatabase(ctrl)

			if tc.prepare != nil {
				tc.prepare(mockDatabase)
			}

			instance := New(mockDatabase)

			err := instance.Delete(tc.uuid)

			tc.expectations(t, err)
		})
	}
}

func TestStorage_Update(t *testing.T) {
	t.Parallel()

	contact := model.Contact{
		UUID: "1",
	}

	tests := []struct {
		name         string
		contact      model.Contact
		prepare      func(db *Mockdatabase)
		expectations func(t assert.TestingT, err error)
	}{
		{
			name:    "Failed to read from database",
			contact: contact,
			prepare: func(db *Mockdatabase) {
				db.EXPECT().
					Read().
					Return(nil, assert.AnError)
			},
			expectations: func(t assert.TestingT, err error) {
				assert.Error(t, err)
			},
		},
		{
			name:    "Not found",
			contact: contact,
			prepare: func(db *Mockdatabase) {
				db.EXPECT().
					Read().
					Return(map[string]Contact{
						"2": {
							UUID: "2",
						},
					}, nil)
			},
			expectations: func(t assert.TestingT, err error) {
				assert.ErrorIs(t, err, model.ErrNotFound)
			},
		},
		{
			name:    "Failed to save",
			contact: contact,
			prepare: func(db *Mockdatabase) {
				db.EXPECT().
					Read().
					Return(map[string]Contact{
						"1": {
							UUID: "1",
						},
					}, nil)

				db.EXPECT().
					Save(map[string]Contact{
						"1": {
							UUID:  "1",
							Links: map[string]string{},
						},
					}).
					Return(assert.AnError)
			},
			expectations: func(t assert.TestingT, err error) {
				assert.Error(t, err)
			},
		},
		{
			name:    "Success",
			contact: contact,
			prepare: func(db *Mockdatabase) {
				db.EXPECT().
					Read().
					Return(map[string]Contact{
						"1": {
							UUID: "1",
						},
					}, nil)

				db.EXPECT().
					Save(map[string]Contact{
						"1": {
							UUID:  "1",
							Links: map[string]string{},
						},
					}).
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
			mockDatabase := NewMockdatabase(ctrl)

			if tc.prepare != nil {
				tc.prepare(mockDatabase)
			}

			instance := New(mockDatabase)

			err := instance.Update(tc.contact)

			tc.expectations(t, err)
		})
	}
}

func TestStorage_Create(t *testing.T) {
	t.Parallel()

	contact := model.Contact{
		UUID:     "1",
		Surname:  "Ершов",
		Name:     "Виталий",
		Birthday: time.Date(2001, 1, 10, 0, 0, 0, 0, time.UTC),
		Phone:    model.NewPhoneFromInt64(79151596781),
		Email:    "vaershov@avito.ru",
		Links: map[model.ContactLink]string{
			model.ContactLinkVk: "vk.com",
		},
	}

	tests := []struct {
		name         string
		contact      model.Contact
		prepare      func(db *Mockdatabase)
		expectations func(t assert.TestingT, err error)
	}{
		{
			name:    "Failed to read from database",
			contact: contact,
			prepare: func(db *Mockdatabase) {
				db.EXPECT().
					Read().
					Return(nil, assert.AnError)
			},
			expectations: func(t assert.TestingT, err error) {
				assert.Error(t, err)
			},
		},
		{
			name:    "Contact already exists",
			contact: contact,
			prepare: func(db *Mockdatabase) {
				db.EXPECT().
					Read().
					Return(map[string]Contact{
						"1": {
							UUID: "1",
						},
					}, nil)
			},
			expectations: func(t assert.TestingT, err error) {
				assert.ErrorIs(t, err, model.ErrAlreadyExists)
			},
		},
		{
			name:    "Failed to save to database",
			contact: contact,
			prepare: func(db *Mockdatabase) {
				db.EXPECT().
					Read().
					Return(map[string]Contact{
						"2": {
							UUID: "2",
						},
					}, nil)

				db.EXPECT().
					Save(map[string]Contact{
						"2": {
							UUID: "2",
						},
						"1": {
							UUID:     "1",
							Surname:  "Ершов",
							Name:     "Виталий",
							Birthday: time.Date(2001, 1, 10, 0, 0, 0, 0, time.UTC),
							Phone:    79151596781,
							Email:    "vaershov@avito.ru",
							Links: map[string]string{
								model.ContactLinkVk: "vk.com",
							},
						},
					}).
					Return(assert.AnError)
			},
			expectations: func(t assert.TestingT, err error) {
				assert.Error(t, err)
			},
		},
		{
			name:    "Success",
			contact: contact,
			prepare: func(db *Mockdatabase) {
				db.EXPECT().
					Read().
					Return(map[string]Contact{
						"2": {
							UUID: "2",
						},
					}, nil)

				db.EXPECT().
					Save(map[string]Contact{
						"2": {
							UUID: "2",
						},
						"1": {
							UUID:     "1",
							Surname:  "Ершов",
							Name:     "Виталий",
							Birthday: time.Date(2001, 1, 10, 0, 0, 0, 0, time.UTC),
							Phone:    79151596781,
							Email:    "vaershov@avito.ru",
							Links: map[string]string{
								model.ContactLinkVk: "vk.com",
							},
						},
					}).
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
			mockDatabase := NewMockdatabase(ctrl)

			if tc.prepare != nil {
				tc.prepare(mockDatabase)
			}

			instance := New(mockDatabase)

			err := instance.Create(tc.contact)

			tc.expectations(t, err)
		})
	}
}
