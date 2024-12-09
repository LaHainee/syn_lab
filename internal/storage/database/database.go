package database

import (
	"encoding/json"
	"fmt"
	"os"

	"contacts/internal/storage"
)

type Database struct {
	path string
}

func New(path string) *Database {
	return &Database{
		path: path,
	}
}

func (d *Database) Save(contacts map[string]storage.Contact) error {
	b, err := json.Marshal(contacts)
	if err != nil {
		return fmt.Errorf("marshall: %w", err)
	}

	err = os.WriteFile(d.path, b, 0644)
	if err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}

func (d *Database) Read() (map[string]storage.Contact, error) {
	b, err := os.ReadFile(d.path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	var dto map[string]storage.Contact
	err = json.Unmarshal(b, &dto)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	return dto, nil
}
