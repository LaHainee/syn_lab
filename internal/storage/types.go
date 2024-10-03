package storage

import (
	"contacts/internal/model"
	"time"
)

type Contact struct {
	UUID     string            `json:"uuid"`
	Surname  string            `json:"surname"`
	Name     string            `json:"name"`
	Birthday time.Time         `json:"birthday"`
	Phone    int64             `json:"phone"`
	Email    string            `json:"email"`
	Links    map[string]string `json:"links"`
}

func dtoToModel(contactDto Contact) model.Contact {
	links := make(map[model.ContactLink]string, len(contactDto.Links))
	for link, value := range contactDto.Links {
		links[model.ContactLink(link)] = value
	}

	return model.Contact{
		UUID:     contactDto.UUID,
		Surname:  contactDto.Surname,
		Name:     contactDto.Name,
		Birthday: contactDto.Birthday,
		Phone:    model.NewPhoneFromInt64(contactDto.Phone),
		Email:    contactDto.Email,
		Links:    links,
	}
}

func modelToDto(contact model.Contact) Contact {
	linksDto := make(map[string]string, len(contact.Links))
	for link, value := range contact.Links {
		linksDto[string(link)] = value
	}

	return Contact{
		UUID:     contact.UUID,
		Surname:  contact.Surname,
		Name:     contact.Name,
		Birthday: contact.Birthday,
		Phone:    contact.Phone.Number(),
		Email:    contact.Email,
		Links:    linksDto,
	}
}
