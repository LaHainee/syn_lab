package main

import (
	"contacts/internal/model"
	"contacts/internal/storage"
	"github.com/google/uuid"
	"time"
)

var (
	birthday = time.Date(2001, 1, 10, 0, 0, 0, 0, time.UTC)
	email    = "fake@gmail.com"
	links    = map[model.ContactLink]string{
		model.ContactLinkVk: "vk.com/vaershov",
	}
)

func main() {
	phone, err := model.NewPhone("+7 (915) 159-67-81")
	if err != nil {
		panic(err)
	}

	storageInstance := storage.New("internal/database/database.json")

	err = create(storageInstance, phone)
	if err != nil {
		panic(err)
	}

	//err := update(storageInstance)
	//if err != nil {
	//	panic(err)
	//}
	//
	//err := deleteContact(storageInstance)
	//if err != nil {
	//	panic(err)
	//}

	//contacts, err := fetch(storageInstance)
	//if err != nil {
	//	panic(err)
	//}
	//
	//for _, contact := range contacts {
	//	fmt.Println(contact)
	//}

	//contacts, err := search(storageInstance, "виталий ершов")
	//if err != nil {
	//	panic(err)
	//}
	//
	//for _, contact := range contacts {
	//	fmt.Println(contact)
	//}
}

func search(storageInstance *storage.Storage, query string) ([]model.Contact, error) {
	return storageInstance.Search(model.SearchRequest{Query: query})
}

func fetch(storageInstance *storage.Storage) ([]model.Contact, error) {
	return storageInstance.Fetch()
}

func deleteContact(storageInstance *storage.Storage) error {
	return storageInstance.Delete("695bb135-63ab-4b2a-bf3c-82339d395e90")
}

func update(storageInstance *storage.Storage, phone model.Phone) error {
	return storageInstance.Update(
		model.Contact{
			UUID:     "695bb135-63ab-4b2a-bf3c-82339d395e90",
			Surname:  "Обновленный",
			Name:     "Семен",
			Birthday: birthday,
			Phone:    phone,
			Email:    email,
			Links:    links,
		},
	)
}

func create(storageInstance *storage.Storage, phone model.Phone) error {
	err := storageInstance.Create(model.Contact{
		UUID:     uuid.NewString(),
		Surname:  "Жихарев",
		Name:     "Семен",
		Birthday: birthday,
		Phone:    phone,
		Email:    email,
		Links:    links,
	})
	if err != nil {
		return err
	}

	err = storageInstance.Create(model.Contact{
		UUID:     uuid.NewString(),
		Surname:  "Ершов",
		Name:     "Виталий",
		Birthday: birthday,
		Phone:    phone,
		Email:    email,
		Links:    links,
	})
	if err != nil {
		return err
	}

	err = storageInstance.Create(model.Contact{
		UUID:     uuid.NewString(),
		Surname:  "Варин",
		Name:     "Дмитрий",
		Birthday: birthday,
		Phone:    phone,
		Email:    email,
		Links:    links,
	})
	if err != nil {
		return err
	}

	return nil
}
