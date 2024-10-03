package storage

import (
	"contacts/internal/model"
	"encoding/json"
	"fmt"
	"golang.org/x/exp/maps"
	"os"
	"slices"
	"strings"
)

type Storage struct {
	db string
}

func New(db string) *Storage {
	return &Storage{
		db: db,
	}
}

// Search – поиск контактов, которые соответствуют запросу
func (s *Storage) Search(request model.SearchRequest) ([]model.Contact, error) {
	contactsDto, err := s.readFromJson()
	if err != nil {
		return nil, err
	}

	words := strings.Fields(request.Query)
	for i, word := range words {
		words[i] = strings.ToLower(word)
	}

	contactsFiltered := make(map[string]model.Contact)

	for _, contactDto := range contactsDto {
		if slices.Contains(words, strings.ToLower(contactDto.Name)) || slices.Contains(words, strings.ToLower(contactDto.Surname)) {
			contactsFiltered[contactDto.UUID] = dtoToModel(contactDto)
		}
	}

	return maps.Values(contactsFiltered), nil
}

// Fetch – получить список всех контактов
func (s *Storage) Fetch() ([]model.Contact, error) {
	contactsDto, err := s.readFromJson()
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
	contactsDto, err := s.readFromJson()
	if err != nil {
		return err
	}

	delete(contactsDto, uuid)

	return s.saveToJson(contactsDto)
}

// Update – обновить контакт, находим контакт по id и перезаписываем его в хранилище
func (s *Storage) Update(contact model.Contact) error {
	contactsDto, err := s.readFromJson()
	if err != nil {
		return err
	}

	_, ok := contactsDto[contact.UUID]
	if !ok {
		return model.ErrNotFound
	}

	contactsDto[contact.UUID] = modelToDto(contact)

	return s.saveToJson(contactsDto)
}

// Create – создать контакт
func (s *Storage) Create(contact model.Contact) error {
	contactsDto, err := s.readFromJson()
	if err != nil {
		return err
	}

	_, ok := contactsDto[contact.UUID]
	if ok {
		return model.ErrAlreadyExists
	}

	contactsDto[contact.UUID] = modelToDto(contact)

	return s.saveToJson(contactsDto)
}

// saveToJson – вспомогательный метод, который инкапсулирует запись в файл
func (s *Storage) saveToJson(contacts map[string]Contact) error {
	b, err := json.Marshal(contacts)
	if err != nil {
		return fmt.Errorf("marshall: %w", err)
	}

	err = os.WriteFile(s.db, b, 0644)
	if err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}

// readFromJson – вспомогательный метод, который инкапсулирует чтение из файла
func (s *Storage) readFromJson() (map[string]Contact, error) {
	b, err := os.ReadFile(s.db)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	var dto map[string]Contact
	err = json.Unmarshal(b, &dto)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	return dto, nil
}
