//go:generate mockgen -source ${GOFILE} -destination mocks_test.go -package ${GOPACKAGE}_test
package delete

type storage interface {
	Delete(uuid string) error
}
