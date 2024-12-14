package contacts_list

import "contacts/internal/model"

type BySurname []model.Contact

func (a BySurname) Len() int           { return len(a) }
func (a BySurname) Less(i, j int) bool { return a[i].Surname < a[j].Surname }
func (a BySurname) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
