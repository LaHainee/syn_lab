//go:generate mockgen -source ${GOFILE} -destination mocks_test.go -package ${GOPACKAGE}_test
package storage

type database interface {
	Read() (map[string]Contact, error)
	Save(contacts map[string]Contact) error
}
