package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"

	"contacts/internal/model"
	"contacts/internal/storage"
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

func generateRandomDate() time.Time {
	rand.Seed(time.Now().UnixNano())

	start := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	randomUnix := rand.Int63n(end.Unix()-start.Unix()) + start.Unix()
	randomDate := time.Unix(randomUnix, 0).UTC()

	return randomDate
}

func main() {
	contactStorage := storage.New("internal/database/database.json")

	contacts := make([]model.Contact, 0)
	for i := 0; i < amount; i++ {
		contacts = append(contacts, model.Contact{
			UUID:     uuid.NewString(),
			Surname:  getRandom(surnames),
			Name:     getRandom(names),
			Birthday: generateRandomDate(),
			Phone:    model.NewPhoneFromInt64(79151596781),
			Email:    "vaershov@avito.ru",
		})
	}

	for _, contact := range contacts {
		err := contactStorage.Create(contact)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Finished")
}
