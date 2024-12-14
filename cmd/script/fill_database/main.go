package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	contactValidator "contacts/internal/domain/validate/contact"
	createContact "contacts/internal/handler/create"
	"contacts/internal/model"
	"contacts/internal/storage"
	"contacts/internal/storage/database"
	"contacts/util/uuid"
)

const amount = 10

var names = []string{
	"Александр", "Дмитрий", "Максим", "Иван", "Сергей",
	"Егор", "Артем", "Никита", "Андрей", "Михаил",
	"Олег", "Константин", "Роман", "Тимур", "Владимир",
}

var surnames = []string{
	"Иванов", "Петров", "Сидоров", "Кузнецов", "Смирнов",
	"Попов", "Лебедев", "Ковалев", "Новиков", "Зайцев",
	"Морозов", "Федоров", "Соловьев", "Григорьев", "Васильев",
}

func getRandom[T any](items []T) T {
	rand.Seed(time.Now().UnixNano())
	return items[rand.Intn(len(items))]
}

func generateRandomDate() string {
	rand.Seed(time.Now().UnixNano())

	start := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	randomUnix := rand.Int63n(end.Unix()-start.Unix()) + start.Unix()
	randomDate := time.Unix(randomUnix, 0).UTC()

	return randomDate.Format("02.01.2006")
}

func main() {
	contactStorage := storage.New(database.New("internal/database/database.json"))
	validator := contactValidator.New()
	uuidGenerator := uuid.NewGenerator()
	createContactHandler := createContact.NewHandler(contactStorage, uuidGenerator, validator)

	contacts := make([]model.ContactForCreate, 0)
	for i := 0; i < amount; i++ {
		contacts = append(contacts, model.ContactForCreate{
			Surname:  getRandom(surnames),
			Name:     getRandom(names),
			Birthday: generateRandomDate(),
			Phone:    "+7 (915) 159-67-81",
			Email:    "vaershov@avito.ru",
		})
	}

	for _, contact := range contacts {
		_, err := createContactHandler.Create(context.Background(), contact)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Finished")
}
