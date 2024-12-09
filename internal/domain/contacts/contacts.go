package contacts

import "contacts/internal/model"

func AllowedLinks() []model.ContactLink {
	return []model.ContactLink{
		model.ContactLinkVk,
	}
}
