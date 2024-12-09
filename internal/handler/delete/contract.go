package delete

type storage interface {
	Delete(uuid string) error
}
