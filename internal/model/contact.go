package model

import "time"

type ContactLink string

const ContactLinkVk = "vk.com"

type Contact struct {
	UUID     string
	Surname  string
	Name     string
	Birthday time.Time
	Phone    Phone
	Email    string
	Links    map[ContactLink]string
}
