package store

type Store interface {
	Set(id string, value string) error
	Get(id string, clear bool) (plain string, err error)
	Verify(id, answer string, clear bool) (match bool, err error)
}
