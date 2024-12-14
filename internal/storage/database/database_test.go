package database_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"contacts/internal/storage"
	. "contacts/internal/storage/database"
)

// Тест для метода Save
func TestDatabase_Save(t *testing.T) {
	// Создаем временный файл для тестирования
	tempFile, err := os.CreateTemp("", "testdb-*.json")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name()) // Удаляем файл после теста

	db := New(tempFile.Name())

	contacts := map[string]storage.Contact{
		"john": {Name: "John Doe", Phone: 79151596781},
		"jane": {Name: "Jane Doe", Phone: 79151596781},
	}

	// Сохраняем контакты
	err = db.Save(contacts)
	require.NoError(t, err, "Неожиданная ошибка при сохранении контактов")

	// Проверяем, что файл был создан и содержит правильные данные
	data, err := os.ReadFile(tempFile.Name())
	require.NoError(t, err)

	var savedContacts map[string]storage.Contact
	err = json.Unmarshal(data, &savedContacts)
	require.NoError(t, err)

	assert.Equal(t, contacts, savedContacts, "Сохраненные контакты должны совпадать с исходными")
}

func TestDatabase_Read(t *testing.T) {
	// Создаем временный файл для тестирования
	tempFile, err := os.CreateTemp("", "testdb-*.json")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name()) // Удаляем файл после теста

	// Записываем тестовые данные в файл
	contacts := map[string]storage.Contact{
		"john": {
			Name:  "John Doe",
			Phone: 79151596781,
		},
	}

	data, err := json.Marshal(contacts)
	require.NoError(t, err)
	_, err = tempFile.Write(data)
	require.NoError(t, err)

	// Создаем экземпляр базы данных
	db := New(tempFile.Name())

	// Читаем контакты
	readContacts, err := db.Read()
	require.NoError(t, err, "Неожиданная ошибка при чтении контактов")

	// Проверяем, что прочитанные контакты совпадают с исходными
	assert.Equal(t, contacts, readContacts, "Прочитанные контакты должны совпадать с исходными")
}

// Тест для обработки ошибки чтения файла
func TestDatabase_Read_FileError(t *testing.T) {
	db := New("non_existent_file.json")

	// Пытаемся читать из несуществующего файла
	contacts, err := db.Read()

	assert.Nil(t, contacts, "Контакты должны быть nil при ошибке чтения")
	assert.Error(t, err, "Ошибка чтения файла не должна быть nil")
}

// Тест для обработки ошибки десериализации
func TestDatabase_Read_UnmarshalError(t *testing.T) {
	// Создаем временный файл для тестирования
	tempFile, err := os.CreateTemp("", "testdb-*.json")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name()) // Удаляем файл после теста

	// Записываем некорректные данные в файл
	_, err = tempFile.Write([]byte("invalid json"))
	require.NoError(t, err)

	// Создаем экземпляр базы данных
	db := New(tempFile.Name())

	// Пытаемся читать из файла с некорректным JSON
	contacts, err := db.Read()

	assert.Nil(t, contacts, "Контакты должны быть nil при ошибке десериализации")
	assert.Error(t, err, "Ошибка десериализации не должна быть nil")
}
