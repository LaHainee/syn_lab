package storage

import (
	"strings"

	"golang.org/x/exp/maps"

	"contacts/internal/model"
)

type Storage struct {
	db database
}

func New(db database) *Storage {
	return &Storage{
		db: db,
	}
}

// Search – поиск контактов, которые соответствуют запросу
func (s *Storage) Search(request model.SearchRequest) ([]model.Contact, error) {
	contactsDto, err := s.db.Read()
	if err != nil {
		return nil, err
	}

	filtered := make(map[string]model.Contact)

	if len(request.Query) == 0 {
		// Для пустого поискового запроса возвращаем все контакты
		for _, contactDto := range contactsDto {
			filtered[contactDto.UUID] = dtoToModel(contactDto)
		}

		return maps.Values(filtered), nil
	}

	words := strings.Fields(request.Query)
	for i, word := range words {
		words[i] = strings.ToLower(word)
	}

	for _, contactDto := range contactsDto {
		for _, word := range words {
			if strings.Contains(strings.ToLower(contactDto.Name), word) || strings.Contains(strings.ToLower(contactDto.Surname), word) {
				filtered[contactDto.UUID] = dtoToModel(contactDto)
			}
		}
	}

	return maps.Values(filtered), nil
}

func (s *Storage) FetchByUuid(uuid string) (model.Contact, error) {
	contactsDto, err := s.db.Read()
	if err != nil {
		return model.Contact{}, err
	}

	for _, contactDto := range contactsDto {
		if contactDto.UUID != uuid {
			continue
		}

		return dtoToModel(contactDto), nil
	}

	return model.Contact{}, model.ErrNotFound
}

// Fetch – получить список всех контактов
func (s *Storage) Fetch() ([]model.Contact, error) {
	contactsDto, err := s.db.Read()
	if err != nil {
		return nil, err
	}

	contacts := make([]model.Contact, 0, len(contactsDto))
	for _, contactDto := range contactsDto {
		contacts = append(contacts, dtoToModel(contactDto))
	}

	return contacts, nil
}

// Delete - удалить контакт по id
func (s *Storage) Delete(uuid string) error {
	contactsDto, err := s.db.Read()
	if err != nil {
		return err
	}

	delete(contactsDto, uuid)

	return s.db.Save(contactsDto)
}

// Update – обновить контакт, находим контакт по id и перезаписываем его в хранилище
func (s *Storage) Update(contact model.Contact) error {
	contactsDto, err := s.db.Read()
	if err != nil {
		return err
	}

	_, ok := contactsDto[contact.UUID]
	if !ok {
		return model.ErrNotFound
	}

	contactsDto[contact.UUID] = modelToDto(contact)

	return s.db.Save(contactsDto)
}

// Create – создать контакт
func (s *Storage) Create(contact model.Contact) error {
	contactsDto, err := s.db.Read()
	if err != nil {
		return err
	}

	_, ok := contactsDto[contact.UUID]
	if ok {
		return model.ErrAlreadyExists
	}

	contactsDto[contact.UUID] = modelToDto(contact)

	return s.db.Save(contactsDto)
}
